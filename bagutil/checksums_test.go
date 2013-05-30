// checksums_test
package bagutil

import (
	//"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var testAlgos []string = []string{"sha1", "sha256", "md5"}

func TestFileChecksum(t *testing.T) {
	testMap := map[string]string{
		"sha1": "da909ba395016f2a64b04d706520db6afa74fc95",
		"md5":  "92d7a9f0f4a30ca782dcae5fe83ca7eb",
	}
	testFile, _ := ioutil.TempFile("", "_GO_")
	testFile.WriteString("Test the checksum")
	testFile.Close()
	for key, sum := range testMap {
		actual := FileChecksum(testFile.Name(), key)
		if sum != actual {
			t.Error("Expected", sum, "but returned", actual, "when checking", key)
		}
	}
	os.Remove(testFile.Name())
}

func TestNewChecksumHash(t *testing.T) {
	for algo := range testAlgos {
		hash, _ := NewChecksumHash(testAlgos[algo])
		if hash == nil {
			t.Error("Expecting a return for", testAlgos[algo], "but returned nil!")
		}
	}
	nilHash, err := NewChecksumHash("badname")
	if nilHash != nil || err == nil {
		t.Error("Expected badhash name to return nil and raise err but it did not!")
	}
}
