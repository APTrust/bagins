package bagutil

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"os"
)

// Takes a filepath as a string and produces a checksum.
func Sha1Checksum(filepath string, algo string) string {
	hsh, _ := newHash(algo)
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
	switch {
	case algo == "md5":
		return md5.New(), nil
	case algo == "sha1":
		return sha1.New(), nil
	}
	return nil, errors.New("Unsupported Hash type.")
}
