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
		return "", err
	}
	return fmt.Sprintf("%x", hsh.Sum(nil)), nil
}

// Utility method to return the operation system seperator as a string.
func PathSeparator() string {
	return string(byte(os.PathSeparator))
}

// Returns a new hash function based on a lookup of the algo string
// passed to the function.  Returns an error if the algo string does not match
// any of the available cryto hashes.
func LookupHash(algo string) (func() hash.Hash, error) {
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
