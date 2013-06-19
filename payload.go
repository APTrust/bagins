package bagins

import (
	"fmt"
	"hash"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Payload struct {
	dir string // Path of the payload directory to manage.
}

// Returns a new Payload struct managing the path provied.
func NewPayload(location string) (*Payload, error) {
	if _, err := os.Stat(path.Clean(location)); os.IsNotExist(err) {
		return nil, fmt.Errorf("Payload directory does not exist! Returned: %v", err)
	}
	p := new(Payload)
	p.dir = path.Clean(location)
	return p, nil
}

func (p *Payload) Name() string {
	return p.dir
}

// TODO Update when this signature settles
func (p *Payload) Add(srcPath string, dstPath string, hsh hash.Hash) (string, error) {

	src, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dstFile := path.Join(p.dir, dstPath)
	if err := os.MkdirAll(path.Dir(dstFile), 0777); err != nil {
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
// method.  Returns a map of the filename and its fixity falue and a
// list of errors.
func (p *Payload) AddAll(src string, hsh hash.Hash) (fxs map[string]string, errs []error) {
	fxs = make(map[string]string)
	visit := func(pth string, info os.FileInfo, err error) error {
		hsh.Reset()
		if err != nil {
			errs = append([]error{err})
		}
		if !info.IsDir() {
			dstPath := strings.TrimPrefix(pth, src)
			fx, err := p.Add(pth, dstPath, hsh)
			if err != nil {
				return err
			}
			fxs[dstPath] = fx
		}
		return nil
	}

	c := make(chan error)

	go func() {
		c <- filepath.Walk(src, visit)
	}()
	if err := <-c; err != nil {
		errs = append([]error{err})
	}

	return
}
