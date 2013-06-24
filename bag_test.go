package bagins_test

import (
	"github.com/APTrust/bagins"
	"github.com/APTrust/bagins/bagutil"
	"os"
	"path"
	"testing"
)

func TestNewBag(t *testing.T) {

	// Use this ChecksumAlgorithm for the tests.
	algo := "sha1"
	hsh, _ := bagutil.LookupHashFunc(algo)
	cs := bagutil.NewChecksumAlgorithm(algo, hsh)

	// It should raise an error if the destination dir does not exist.
	badLocation := path.Join(os.TempDir(), "/GOTESTNOT_EXISTs/")
	_, err := bagins.NewBag(badLocation, "_GOFAILBAG_", cs)
	if err == nil {
		t.Error("NewBag function does not recognize when a directory does not exist!")
	}

	// It should raise an error if the bag already exists.
	os.MkdirAll(path.Join(badLocation, "_GOFAILBAG_"), 0766)
	defer os.RemoveAll(badLocation)

	_, err = bagins.NewBag(badLocation, "_GOFAILBAG_", cs)
	if err == nil {
		t.Error("Error not thrown when bag already exists as expected.")
	}

	// It should create a bag without any errors.

	bag, err := bagins.NewBag(os.TempDir(), "_GOTESTBAG_", cs)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(bag.Path())

	// Confirm existance of basic bag structure.
	if _, err = os.Stat(path.Join(os.TempDir(), "_GOTESTBAG_")); os.IsNotExist(err) {
		t.Error("Bag directory does not exist!")
	}
	if _, err = os.Stat(path.Join(os.TempDir(), "_GOTESTBAG_", "data")); os.IsNotExist(err) {
		t.Error("Data directory does not exist!")
	}
	if _, err = os.Stat(path.Join(os.TempDir(), "_GOTESTBAG_", "bagit.txt")); os.IsNotExist(err) {
		t.Error("bagit.txt does not exist!")
	}
}
