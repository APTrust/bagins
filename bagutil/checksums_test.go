// checksums_test
package bagutil

import (
	"io/ioutil"
	"os"
	"testing"
)

var testAlgos []string = []string{"sha1", "sha256", "md5"}

func TestNewCheckAlgorithm(t *testing.T) {
	for idx := range testAlgos {
		hsh, _ := LookupHash(testAlgos[idx])
		chkAlgo := NewCheckAlgorithm(testAlgos[idx], hsh)
		if chkAlgo == nil {
			t.Error("Returned nil for check algorithm")
		}
	}
}

func TestNewCheckByName(t *testing.T) {
	for idx := range testAlgos {
		_, err := NewCheckByName(testAlgos[idx])
		if err != nil {
			t.Error(err)
		}
	}
}

func TestFileChecksum(t *testing.T) {
	testMap := map[string]string{
		"sha1": "da909ba395016f2a64b04d706520db6afa74fc95",
		"md5":  "92d7a9f0f4a30ca782dcae5fe83ca7eb",
	}
	testFile, _ := ioutil.TempFile("", "_GO_")
	testFile.WriteString("Test the checksum")
	testFile.Close()
	for key, sum := range testMap {
		hsh, _ := LookupHash(key)
		actual, err := FileChecksum(testFile.Name(), hsh)
		if err != nil {
			t.Error(err)
		}
		if sum != actual {
			t.Error("Expected", sum, "but returned", actual, "when checking", key)
		}
	}
	os.Remove(testFile.Name())
}

func TestLookupHash(t *testing.T) {
	for algo := range testAlgos {
		hash, _ := LookupHash(testAlgos[algo])
		if hash == nil {
			t.Error("Expecting a return for", testAlgos[algo], "but returned nil!")
		}
	}
	nilHash, err := LookupHash("badname")
	if nilHash != nil || err == nil {
		t.Error("Expected badhash name to return nil and raise err but it did not!")
	}
}
