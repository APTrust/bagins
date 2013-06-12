package bagins

import (
	"fmt"
	"github.com/APTrust/bagins/bagutil"
	"hash"
	"os"
	"path"
)

type Payload struct {
	path string // Path of the payload directory to manage.
}

// Returns a new Payload struct managing the path provied.
func NewPayload(location string) (*Payload, error) {
	p := new(Payload)
	if _, err := os.Stat(path.Clean(location)); os.IsNotExist(err) {
		return nil, fmt.Errorf("Payload directory does not exist! Returned: %v", err)
	}
	p.path = path.Clean(location)
	return p, nil
}

// TODO Update when this signature settles
func (p *Payload) Add(srcPath string, dstPath string, hsh hash.Hash) (string, error) {
	chkSum, err := bagutil.FileChecksum(filepath, hsh)
	if err != nil {
		return "", err
	}
	src, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(path.Join(p.path, dstPath))
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	wrtn, err := io.Copy(dst, src)
	return wrtn, err
}

// Performs an add on every file under the directory supplied to the
// method.  Returns a map of the filename and its fixity falue and a
// list of errors.
func (p *Payload) AddAll(dir string, hsh hash.Hash) (fxs map[string]string, errs []error) {
	return
}
