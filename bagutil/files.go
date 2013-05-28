package bagutil

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
)

// Test the output of some hashes to assess use in fixity checks.
func Sha1Checksum(file *os.File) string {
	hsh := sha1.New()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	hsh.Write(fileBytes)
	byteSum := hsh.Sum(nil)
	return fmt.Sprintf("%x", byteSum) // Convert to base16 on formatting.
}
