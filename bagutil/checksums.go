package bagutil

import (
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

type HashMaker func() hash.Hash

type CheckAlgorithm struct {
	Name string
	Hash hash.Hash
}

func NewCheckAlgorithm(name string, hsh hash.Hash) *CheckAlgorithm {
	h := new(CheckAlgorithm)
	h.Name = name
	h.Hash = hsh
	return h
}

// Convienence method that looks up a checksum by name and assigns it
// properly or returns an error.
func NewCheckByName(name string) (*CheckAlgorithm, error) {
	hsh, err := LookupHash(name)
	if err != nil {
		return nil, err
	}
	h := new(CheckAlgorithm)
	h.Name = strings.ToLower(name)
	h.Hash = hsh
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
	return nil, errors.New("Unsupported hash value.")
}
