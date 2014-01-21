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
// a checksum value as calulated by the provided hash.  Returns the checksum
// string and any error encountered
func (p *Payload) Add(srcPath string, dstPath string, hsh hash.Hash) (string, error) {

	src, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	// defer src.Close()

	// TODO simplify this! returns on windows paths are messing with me so I'm
	// going through this step wise.
	dstFile := filepath.Join(p.dir, dstPath)
	if err := os.MkdirAll(filepath.Dir(dstFile), 0766); err != nil {
		return "", err
	}

	dst, err := os.Create(dstFile)
	if err != nil {
		return "", err
	}
	//defer dst.Close()

	wrtr := io.MultiWriter(dst, hsh)

	_, err = io.Copy(wrtr, src)
	if err != nil {
		return "", err
	}
	chkSum := fmt.Sprintf("%x", hsh.Sum(nil))
	src.Close()
	dst.Close()
	return chkSum, err
}

// Performs an add on every file under the directory supplied to the
// method.  Returns a map of the filenames and its fixity value based
// on the hash function passed and a slice of errors if there were any.
func (p *Payload) AddAll(src string, hsh hash.Hash) (fxs map[string]string, errs []error) {

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
	queue := make(chan bool, 100)
	wg := sync.WaitGroup{}
	for idx := range files {
		queue <- true
		wg.Add(1)
		go func(file string, src string, hsh hash.Hash) {
			dstPath := strings.TrimPrefix(file, src)
			fx, err := p.Add(file, dstPath, hsh)
			if err != nil {
				errs = append(errs, err)
			}
			fxs[dstPath] = fx
			<-queue
			wg.Done()
		}(files[idx], src, hsh)
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
