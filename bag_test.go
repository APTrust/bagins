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

	// It should find all of the following files and directories.
	if _, err = os.Stat(path.Join(os.TempDir(), "_GOTESTBAG_")); os.IsNotExist(err) {
		t.Error("Bag directory does not exist!")
	}
	if data, err := os.Stat(path.Join(os.TempDir(), "_GOTESTBAG_", "data")); os.IsNotExist(err) || !data.IsDir() {
		t.Error("Data directory does not exist or is not a directory!")
	}
	if _, err = os.Stat(path.Join(bag.Path(), "bagit.txt")); os.IsNotExist(err) {
		bi, err := bag.BagInfo()
		if err != nil {
			t.Error(err)
		}
		t.Errorf("bagit.txt does not exist! %s", bi.Name())
	}
	if _, err = os.Stat(path.Join(os.TempDir(), "_GOTESTBAG_", "manifest-sha1.txt")); os.IsNotExist(err) {
		t.Error("manifest-sha1.txt does not exist!")
	}
}

// Test closing the bag in various secenarios for expected behavior.
// func TestBagClose(t *testing.T) {
// 	// Use this ChecksumAlgorithm for the tests.
// 	algo := "sha1"
// 	hsh, _ := bagutil.LookupHashFunc(algo)
// 	cs := bagutil.NewChecksumAlgorithm(algo, hsh)

// 	bag, err := bagins.NewBag(os.TempDir(), "_GOTESTBAG_", cs)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer os.RemoveAll(bag.Path())

// 	// It should find all of the following files and directories.
// 	if _, err = os.Stat(path.Join(os.TempDir(), "_GOTESTBAG_")); os.IsNotExist(err) {
// 		t.Error("Bag directory does not exist!")
// 	}
// 	if data, err := os.Stat(path.Join(os.TempDir(), "_GOTESTBAG_", "data")); os.IsNotExist(err) || !data.IsDir() {
// 		t.Error("Data directory does not exist or is not a directory!")
// 	}
// 	if _, err = os.Stat(path.Join(os.TempDir(), "_GOTESTBAG_", "bagit.txt")); os.IsNotExist(err) {
// 		t.Error("bagit.txt does not exist!")
// 	}
// 	if _, err = os.Stat(path.Join(os.TempDir(), "_GOTESTBAG_", "manifest-sha1.txt")); os.IsNotExist(err) {
// 		t.Error("manifest-sha1.txt does not exist!")
// 	}

// }
