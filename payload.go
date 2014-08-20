package bagins

/*

"Faithless is he that says farewell when the road darkens."

- Gimli

*/

import (
	"fmt"
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
// a checksum value as calulated by the provided hash.  Returns the checksum
// string and any error encountered
func (p *Payload) Add(srcPath string, dstPath string, m *Manifest) (string, error) {

	hsh := m.hashFunc()
	src, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dstFile := filepath.Join(p.dir, dstPath)

	var wrtr io.Writer = nil

	absSrcPath, err := filepath.Abs(srcPath)
	if err != nil {
		return "", err
	}
	absDestPath, err := filepath.Abs(dstFile)
	if err != nil {
		return "", err
	}

	// If src and dst are the same, copying with destroy the src.
	// Just compute the hash.
	if absSrcPath == absDestPath {
		wrtr = io.MultiWriter(hsh)
	} else {
		// TODO simplify this! returns on windows paths are messing with me so I'm
		// going through this step wise.
		if err := os.MkdirAll(filepath.Dir(dstFile), 0766); err != nil {
			return "", err
		}
		dst, err := os.Create(dstFile)
		if err != nil {
			return "", err
		}
		wrtr = io.MultiWriter(dst, hsh)
		defer dst.Close()
	}

	_, err = io.Copy(wrtr, src)
	if err != nil {
		return "", err
	}
	chkSum := fmt.Sprintf("%x", hsh.Sum(nil))
	return chkSum, err
}

// Performs an add on every file under the directory supplied to the
// method.  Returns a map of the filenames and its fixity value based
// on the hash function passed and a slice of errors if there were any.
func (p *Payload) AddAll(src string, m *Manifest) (fxs map[string]string, errs []error) {

	fxs = make(map[string]string)

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
	for idx := range files {
		queue <- true
		wg.Add(1)
		go func(file string, src string, m *Manifest) {
			dstPath := strings.TrimPrefix(file, src)
			fx, err := p.Add(file, dstPath, m)
			if err != nil {
				errs = append(errs, err)
			}
			fxs[dstPath] = fx
			<-queue
			wg.Done()
		}(files[idx], src, m)
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
