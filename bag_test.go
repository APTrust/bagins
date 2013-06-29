package bagins_test

import (
	"github.com/APTrust/bagins"
	"github.com/APTrust/bagins/bagutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewBag(t *testing.T) {

	// Use this ChecksumAlgorithm for the tests.
	algo := "sha1"
	hsh, _ := bagutil.LookupHashFunc(algo)
	cs := bagutil.NewChecksumAlgorithm(algo, hsh)

	// It should raise an error if the destination dir does not exist.
	badLocation := filepath.Join(os.TempDir(), "/GOTESTNOT_EXISTs/")
	_, err := bagins.NewBag(badLocation, "_GOFAILBAG_", cs)
	if err == nil {
		t.Error("NewBag function does not recognize when a directory does not exist!")
	}

	// It should raise an error if the bag already exists.
	os.MkdirAll(filepath.Join(badLocation, "_GOFAILBAG_"), 0766)
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
	if _, err = os.Stat(filepath.Join(os.TempDir(), "_GOTESTBAG_")); os.IsNotExist(err) {
		t.Error("Bag directory does not exist!")
	}
	if data, err := os.Stat(filepath.Join(os.TempDir(), "_GOTESTBAG_", "data")); os.IsNotExist(err) || !data.IsDir() {
		t.Error("Data directory does not exist or is not a directory!")
	}
	if _, err = os.Stat(filepath.Join(bag.Path(), "bagit.txt")); os.IsNotExist(err) {
		bi, err := bag.BagInfo()
		if err != nil {
			t.Error(err)
		}
		t.Errorf("bagit.txt does not exist! %s", bi.Name())
	}
	if _, err = os.Stat(filepath.Join(os.TempDir(), "_GOTESTBAG_", "manifest-sha1.txt")); os.IsNotExist(err) {
		t.Error("manifest-sha1.txt does not exist!")
	}
}

// It should place an appropriate file in the data directory and add the fixity to the manifest.
func TestAddFile(t *testing.T) {
	// Setup the test file to add for the test.
	fi, _ := ioutil.TempFile("", "TEST_GO_ADDFILE_")
	fi.WriteString("Test the checksum")
	fi.Close()
	defer os.Remove(fi.Name())

	// Setup the Test Bag
	algo := "sha1"
	hsh, _ := bagutil.LookupHashFunc(algo)
	cs := bagutil.NewChecksumAlgorithm(algo, hsh)

	bag, err := bagins.NewBag(os.TempDir(), "_GOTESTBAG_", cs)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(bag.Path())

	// It should return an error when trying to add a file that doesn't exist.
	if err = bag.AddFile("idontexist.txt", "idontexist.txt"); err == nil {
		t.Errorf("Adding a nonexistant file did not generate an error!")
	}

	// It should and a file to the data directory and generate a fixity value.
	expFile := "my/nested/dir/mytestfile.txt"
	if err = bag.AddFile(fi.Name(), expFile); err != nil {
		t.Error(err)
	}

	// It should have created the file in the payload directory.
	_, err = os.Stat(filepath.Join(bag.Path(), "data", expFile))
	if err != nil {
		t.Error("Testing if payload file created:", err)
	}

	// It should have calulated the fixity and put it in the manifest.
	mf, _ := bag.Manifest()
	expKey := filepath.Join("data", expFile)
	fx, ok := mf.Data[expKey]
	if !ok {
		t.Error("Unable to find entry in manfest: ", expKey)
	}
	if len(fx) != 40 {
		t.Errorf("Expected %d character fixity but returned: %d", 32, len(fx))
	}
}

func TestAddDir(t *testing.T) {

	// Setup source files to test
	srcDir, _ := ioutil.TempDir("", "_GOTEST_PAYLOAD_SRC_")
	for i := 0; i < 50; i++ {
		fi, _ := ioutil.TempFile(srcDir, "TEST_GO_ADDFILE_")
		fi.WriteString("Test the checksum")
		fi.Close()
	}
	defer os.RemoveAll(srcDir)

	// Setup the test bag
	algo := "sha1"
	hsh, _ := bagutil.LookupHashFunc(algo)
	cs := bagutil.NewChecksumAlgorithm(algo, hsh)

	bag, err := bagins.NewBag(os.TempDir(), "_GOTESTBAG_", cs)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(bag.Path())

	// It should produce no errors
	if errs := bag.AddDir(srcDir); len(errs) != 0 {
		t.Error(errs)
	}

	// It should produce 50 manifest entries

	// It should generate entries in the manifest
	mf, _ := bag.Manifest()

	// It should produce 50 manifest entries
	if len(mf.Data) != 50 {
		t.Error("Expected 50 manifest entries but returned", len(mf.Data))
	}
	// It should contain the proper checksums for each file.
	for key, fx := range mf.Data {
		expFx := "da909ba395016f2a64b04d706520db6afa74fc95"
		expPfx := filepath.Join("data", "TEST_GO_ADDFILE_")

		if fx != expFx {
			t.Error("Fixity error!", fx, "does not match expected", expFx)
		}
		if !strings.HasPrefix(key, expPfx) {
			t.Error(key, "does not start with", expPfx)
		}
	}
}
