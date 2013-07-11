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

type ChecksumAlgorithm struct {
	name string
	hsh  func() hash.Hash
}

func (cs *ChecksumAlgorithm) New() hash.Hash {
	return cs.hsh()
}

func (cs *ChecksumAlgorithm) Algo() func() hash.Hash {
	return cs.hsh
}

func (cs *ChecksumAlgorithm) Name() string {
	return cs.name
}

func NewChecksumAlgorithm(name string, hsh func() hash.Hash) *ChecksumAlgorithm {
	cs := new(ChecksumAlgorithm)
	cs.name = name
	cs.hsh = hsh
	return cs
}

// Convienence method that looks up a checksum by name and assigns it
// properly or returns an error.
func NewCheckByName(name string) (*ChecksumAlgorithm, error) {
	hsh, err := LookupHashFunc(name)
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
	src, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	_, err = io.Copy(hsh, src)
	if err != nil {
		panic(err)
	}
	byteSum := hsh.Sum(nil)
	return fmt.Sprintf("%x", byteSum), nil
}

// Returns a new hash.Hash as indicated by the algo string.
// TODO change this to return a func() hash.Hash
func LookupHash(algo string) (hash.Hash, error) {
	hsh, err := LookupHashFunc(algo)
	if err != nil {
		return nil, err
	}
	return hsh(), nil
}

// Returns a new hash function based on a lookup of the algo string
// passed to the function.  Returns an error if the algo string does not match
// any of the available cryto hashes.
func LookupHashFunc(algo string) (func() hash.Hash, error) {
	switch strings.ToLower(algo) {
	case "md5":
		return crypto.MD5.New, nil
	case "sha1":
		return crypto.SHA1.New, nil
	case "sha256":
		return crypto.SHA256.New, nil
	case "sha512":
		return crypto.SHA512.New, nil
	case "sha224":
		return crypto.SHA224.New, nil
	case "sha384":
		return crypto.SHA384.New, nil
	}
	return nil, fmt.Errorf("Invalid hash name %s:  Must be one of md5, sha1, sha256, sha512, sha224, sha284")
}
