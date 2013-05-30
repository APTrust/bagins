package bagutil

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
)

// Performs a checksum on a file located at `filepath` using `algo`
func FileChecksum(filepath string, algo string) string {
	hsh, err := NewChecksumHash(algo)
	if err != nil {
		panic(err)
	}
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(hsh, file)
	if err != nil {
		panic(err)
	}
	byteSum := hsh.Sum(nil)
	return fmt.Sprintf("%x", byteSum) // Convert to base16 on formatting.
}

// Returns a new hash.Hash as indicated by the algo string.
func NewChecksumHash(algo string) (hash.Hash, error) {
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
