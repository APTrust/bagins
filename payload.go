package bagins

import (
	"fmt"
	//"github.com/APTrust/bagins/bagutil"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	defer src.Close()

	dstFile := filepath.Join(p.dir, dstPath)
	if err := os.MkdirAll(filepath.Dir(dstFile), 0777); err != nil {
		return "", err
	}

	dst, err := os.Create(dstFile)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	wrtr := io.MultiWriter(dst, hsh)

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
func (p *Payload) AddAll(src string, hsh func() hash.Hash) (fxs map[string]string, errs []error) {

	fxs = make(map[string]string)

	// Collect files to add in scr directory.
	var files []string
	visit := func(pth string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append([]string{pth})
		}
		return nil
	}

	if err := filepath.Walk(src, visit); err != nil {
		errs = append([]error{err})
	}

	// Perform Payload.Add on each file found in src under a goroutine.
	c := make(chan bool)
	for idx := range files {
		go func() {
			dstPath := strings.TrimPrefix(files[idx], src)
			fx, err := p.Add(files[idx], dstPath, hsh())
			if err != nil {
				errs = append([]error{err})
			}
			fxs[dstPath] = fx
			c <- true
		}()
	}

	// wait for all go routines to reply.
	for i := 0; i < len(files); i++ {
		<-c // Tick off as goroutines return true
	}

	return
}
