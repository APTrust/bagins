package bagutil

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"os"
)

// Takes a filepath as a string and produces a checksum.
func FileChecksum(filepath string, algo string) string {
	hsh, err := newHash(algo)
	if err != nil {
		panic(err)
	}
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	hsh.Write(fileBytes)
	byteSum := hsh.Sum(nil)
	return fmt.Sprintf("%x", byteSum) // Convert to base16 on formatting.
}

func newHash(algo string) (hash.Hash, error) {
	switch algo {
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	case "sha512":
		return sha512.New(), nil
	}
	return nil, errors.New("Unsupported hash value.")
}
