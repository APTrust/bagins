package bagins

import (
	"fmt"
	"github.com/APTrust/bagins/bagutil"
	"hash"
	"io"
	"os"
	"path"
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

// TODO Update when this signature settles
func (p *Payload) Add(srcPath string, dstPath string, hsh hash.Hash) (string, error) {
	chkSum, err := bagutil.FileChecksum(srcPath, hsh)
	if err != nil {
		return "", err
	}
	src, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(path.Join(p.dir, dstPath))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return chkSum, err
}

// Performs an add on every file under the directory supplied to the
// method.  Returns a map of the filename and its fixity falue and a
// list of errors.
func (p *Payload) AddAll(dir string, hsh hash.Hash) (fxs map[string]string, errs []error) {
	return
}
