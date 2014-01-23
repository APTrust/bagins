package bagutil

/*

"Deeds will not be less valiant because they are unpraised."

- Aragorn

*/

import (
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

// Convienence structure that ties a checksum hash to a string indicating its name/type.
type ChecksumAlgorithm struct {
	name string
	hsh  hash.Hash
}

// Resets the hash function for fresh use and returns it.
func (cs *ChecksumAlgorithm) Algo() hash.Hash {
	cs.hsh.Reset()
	return cs.hsh
}

// Returns the name of the checksum algorithm being used as set when created.
func (cs *ChecksumAlgorithm) Name() string {
	return cs.name
}

// Returns a pointer to a new checksum algorithm.
func NewChecksumAlgorithm(name string, hsh hash.Hash) *ChecksumAlgorithm {
	cs := new(ChecksumAlgorithm)
	cs.name = name
	cs.hsh = hsh
	return cs
}

// Convienence method that looks up a checksum by name and assigns it
// properly or returns an error.
func NewCheckByName(name string) (*ChecksumAlgorithm, error) {
	hsh, err := LookupHash(name)
	if err != nil {
		return nil, err
	}
	h := new(ChecksumAlgorithm)
	h.name = strings.ToLower(name)
	h.hsh = hsh
	return h, nil
}

// Performs a checksum with the hsh.Hash.Sum() method passed to the function
// and returns the hex value of the resultant string or an error
func FileChecksum(filepath string, hsh hash.Hash) (string, error) {
	hsh.Reset()
	src, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	_, err = io.Copy(hsh, src)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hsh.Sum(nil)), nil
}

// Returns a new hash function based on a lookup of the algo string
// passed to the function.  Returns an error if the algo string does not match
// any of the available cryto hashes.
func LookupHash(algo string) (hash.Hash, error) {
	switch strings.ToLower(algo) {
	case "md5":
		return crypto.MD5.New(), nil
	case "sha1":
		return crypto.SHA1.New(), nil
	case "sha256":
		return crypto.SHA256.New(), nil
	case "sha512":
		return crypto.SHA512.New(), nil
	case "sha224":
		return crypto.SHA224.New(), nil
	case "sha384":
		return crypto.SHA384.New(), nil
	}
	return nil, fmt.Errorf("Invalid hash name %s:  Must be one of md5, sha1, sha256, sha512, sha224, sha284")
}
