// manifest_test
package bagins_test

import (
	"bagins"
	"io/ioutil"
	"os"
	"testing"
)

func TestRunChecksums(t *testing.T) {
	testFile, _ := ioutil.TempFile("", "_GO_")
	testFile.WriteString("Test the checksum")
	testFile.Close()

	mfst := bagins.Manifest{Algo: "sha1", Data: make(map[string]string)}
	mfst.Data[testFile.Name()] = ""
	mfst.RunChecksums()

	// Checksum for file should now be generated.
	if mfst.Data[testFile.Name()] != "da909ba395016f2a64b04d706520db6afa74fc95" {
		t.Error("File checksum not accurantly generated!")
	}

	// Check that it throws an error if mismatch checksum.
	mfst.Data[testFile.Name()] = "frodo lives!"
	errList := mfst.RunChecksums()
	if len(errList) == 0 {
		t.Error("Invalid Checksums not being detected!")
	}
	os.Remove(testFile.Name()) // Remove the test file.
}
