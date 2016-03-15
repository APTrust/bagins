package bagins

/*

"Faithless is he that says farewell when the road darkens."

- Gimli

*/

import (
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Payloads describes a filepath location to serve as the data directory of
// a Bag and methods around managing content inside of it.
type Payload struct {
	dir string // Path of the payload directory to manage.
}

// Returns a new Payload struct managing the path provied.
func NewPayload(location string) (*Payload, error) {
	if _, err := os.Stat(filepath.Clean(location)); os.IsNotExist(err) {
		return nil, fmt.Errorf("Payload directory does not exist! Returned: %v", err)
	}
	p := new(Payload)
	p.dir = filepath.Clean(location)
	return p, nil
}

func (p *Payload) Name() string {
	return p.dir
}

// Adds the file at srcPath to the payload directory as dstPath and returns
// a checksum value as calulated by the provided hash. This function also
// writes the checksums into the proper manifests, so you don't have to.
//
// Param manifests should be a slice of payload manifests, which you can get
// from a bag by calling:
//
// bag.GetManifests(PayloadManifest)
//
// Returns the checksums in the form of a map whose keys are the algorithms
// and values are the digests.
//
// If you have an md5 manifest and a sha256 manifest, you'll get back a map
// that looks like this:
//
// checksums["md5"] = "0a0a0a0a"
// checksums["sha256"] = "0b0b0b0b"
func (p *Payload) Add(srcPath string, dstPath string, manifests []*Manifest) (map[string]string, error) {

	src, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dstFile := filepath.Join(p.dir, dstPath)

	var wrtr io.Writer = nil

	absSrcPath, err := filepath.Abs(srcPath)
	if err != nil {
		return nil, err
	}
	absDestPath, err := filepath.Abs(dstFile)
	if err != nil {
		return nil, err
	}

	hashWriters := make([]io.Writer, 0)
	hashFunctions := make([]hash.Hash, 0)
	hashFunctionNames := make([]string, 0)

	// Note that we're putting the same object into
	// hashWriters and hashFunctions, because we need
	// them to behave as both io.Writer and hash.Hash.
	for _, m := range manifests {
		hashObj := m.hashFunc()
		hashWriters = append(hashWriters, hashObj)
		hashFunctions = append(hashFunctions, hashObj)
		hashFunctionNames = append(hashFunctionNames, m.Algorithm())
	}


	// If src and dst are the same, copying with destroy the src.
	// Just compute the hash.
	if absSrcPath == absDestPath {
		wrtr = io.MultiWriter(hashWriters...)
	} else {
		// TODO simplify this! returns on windows paths are messing with me so I'm
		// going through this step wise.
		if err := os.MkdirAll(filepath.Dir(dstFile), 0766); err != nil {
			return nil, err
		}
		dst, err := os.Create(dstFile)
		if err != nil {
			return nil, err
		}
		// Append the destination file to our group of hashWriters,
		// so the file actually gets copied.
		hashWriters = append(hashWriters, dst)
		wrtr = io.MultiWriter(hashWriters...)
		defer dst.Close()
	}

	// Copy the file and compute the hashes. Note that if src and dest
	// are the same, we're only only computing the hash without actually
	// copying the bits.
	_, err = io.Copy(wrtr, src)
	if err != nil {
		return nil, err
	}

	// Calculate the checksums in hex format, so we can return them
	// and write them into the manifests.
	checksums := make(map[string]string)
	for index := range hashFunctions {
		manifest := manifests[index]
		name := hashFunctionNames[index]
		hashFunc := hashFunctions[index]
		digest := fmt.Sprintf("%x", hashFunc.Sum(nil))
		checksums[name] = digest

		// Add the path and digest to the manifest
		manifest.Data[filepath.Join("data", dstPath)] = digest
	}
	return checksums, err
}

// Performs an add on every file under the directory supplied to the
// method. Returns a map of filenames and fixity values based
// on the hash function in the manifests.
//
// Param manifests should be a slice of payload manifests, which you can get
// from a bag by calling:
//
// bag.GetManifests(PayloadManifest)
//
// If you have an md5 manifest and a sha256 manifest, you'll get back a map
// that looks like this:
//
// checksums["file1.txt"] = { "md5": "0a0a0a0a", "sha256": "0b0b0b0b" }
// checksums["file2.xml"] = { "md5": "1a1a1a1a", "sha256": "1b1b1b1b" }
// checksums["file3.jpg"] = { "md5": "2a2a2a2a", "sha256": "2b2b2b2b" }
func (p *Payload) AddAll(src string, manifests []*Manifest) (checksums map[string]map[string]string, errs []error) {

	checksums = make(map[string]map[string]string)

	// Collect files to add in scr directory.
	var files []string
	visit := func(pth string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, pth)
		}
		return err
	}

	if err := filepath.Walk(src, visit); err != nil {
		errs = append(errs, err)
	}
	// Perform Payload.Add on each file found in src under a goroutine.
	queue := make(chan bool, 5)
	wg := sync.WaitGroup{}
	for index := range files {
		queue <- true
		wg.Add(1)
		go func(file string, src string, manifests []*Manifest) {
			dstPath := strings.TrimPrefix(file, src)
			fixities, err := p.Add(file, dstPath, manifests)
			if err != nil {
				errs = append(errs, err)
			}
			checksums[dstPath] = fixities
			<-queue
			wg.Done()
		}(files[index], src, manifests)
	}

	wg.Wait()

	return
}

// Returns the octetstream sum and number of files of all the files in the
// payload directory.  See the BagIt specification "Oxsum" field of the
// bag-info.txt file for more information.
func (p *Payload) OctetStreamSum() (int64, int) {
	var sum int64
	var count int

	visit := func(pth string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			sum = sum + info.Size()
			count = count + 1
		}
		return err
	}

	filepath.Walk(p.dir, visit)

	return sum, count
}
