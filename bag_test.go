package bagins_test

import (
	"github.com/APTrust/bagins"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	FIXSTRING string = "The quick brown fox jumps over the lazy dog."
	FIXVALUE  string = "e4d909c290d0fb1ca068ffaddf22cbd0"
)

func setupTestBag(bagName string) (*bagins.Bag, error) {
	bag, err := bagins.NewBag(os.TempDir(), bagName, "md5")
	if err != nil {
		return nil, err
	}
	return bag, nil
}

func TestNewBag(t *testing.T) {

	// It should raise an error if the destination dir does not exist.
	badLocation := filepath.Join(os.TempDir(), "/GOTESTNOT_EXISTs/")
	_, err := bagins.NewBag(badLocation, "_GOFAILBAG_", "md5")
	if err == nil {
		t.Error("NewBag function does not recognize when a directory does not exist!")
	}

	// It should raise an error if the bag already exists.
	os.MkdirAll(filepath.Join(badLocation, "_GOFAILBAG_"), 0766)
	defer os.RemoveAll(badLocation)

	_, err = bagins.NewBag(badLocation, "_GOFAILBAG_", "md5")
	if err == nil {
		t.Error("Error not thrown when bag already exists as expected.")
	}

	// It should create a bag without any errors.
	bagName := "_GOTEST_NEWBAG_"
	bag, err := setupTestBag("_GOTEST_NEWBAG_")
	defer os.RemoveAll(bag.Path())

	// It should find all of the following files and directories.
	if _, err = os.Stat(filepath.Join(os.TempDir(), bagName)); os.IsNotExist(err) {
		t.Error("Bag directory does not exist!")
	}
	if data, err := os.Stat(filepath.Join(os.TempDir(), bagName, "data")); os.IsNotExist(err) || !data.IsDir() {
		t.Error("Data directory does not exist or is not a directory!")
	}
	if _, err = os.Stat(filepath.Join(bag.Path(), "bagit.txt")); os.IsNotExist(err) {
		bi, err := bag.BagInfo()
		if err != nil {
			t.Error(err)
		}
		t.Errorf("bagit.txt does not exist! %s", bi.Name())
	}
	if _, err = os.Stat(filepath.Join(os.TempDir(), bagName, "manifest-md5.txt")); os.IsNotExist(err) {
		t.Error("manifest-md5.txt does not exist!")
	}
}

func TestReadBag(t *testing.T) {

	// It should return an error when passed a path that doesn't exist.
	badPath := "/thispath/isbad"
	if _, err := bagins.ReadBag(badPath, []string{}, ""); err == nil {
		t.Errorf("Path %s not detected as bad as expected.", badPath)
	}

	// It should return an error if it isn't passed a path to a directory.
	fi, _ := ioutil.TempFile("", "TEST_GO_READBAG_")
	fi.WriteString("Test file please delete.")
	fi.Close()
	defer os.Remove(fi.Name())

	if _, err := bagins.ReadBag(fi.Name(), []string{}, ""); err == nil {
		t.Error("Readbag should thrown an error when trying to open a file: %s", fi.Name())
	}

	// It should return an error if the directory does not contain a data subdirectory.
	pDir, _ := ioutil.TempDir("", "_GOTEST_Payload_")
	defer os.Remove(pDir)

	if _, err := bagins.ReadBag(pDir, []string{}, ""); err == nil {
		t.Errorf("Not returning expected error when directory has no data subdirectory for %s", pDir)
	}

	// It should return an error if it can't determine a valid manifest file.

	// It should return an error if it return a valid bag-info.txt

	// It should return a valid bag.

}

// It should place an appropriate file in the data directory and add the fixity to the manifest.
func TestAddFile(t *testing.T) {
	// Setup the test file to add for the test.
	fi, _ := ioutil.TempFile("", "TEST_GO_ADDFILE_")
	fi.WriteString("Test the checksum")
	fi.Close()
	defer os.Remove(fi.Name())

	// Setup the Test Bag
	bag, _ := setupTestBag("_GOTEST_BAG_ADDFILE_")
	defer os.RemoveAll(bag.Path())

	// It should return an error when trying to add a file that doesn't exist.
	if err := bag.AddFile("idontexist.txt", "idontexist.txt"); err == nil {
		t.Errorf("Adding a nonexistant file did not generate an error!")
	}

	// It should add a file to the data directory and generate a fixity value.
	expFile := "my/nested/dir/mytestfile.txt"
	if err := bag.AddFile(fi.Name(), expFile); err != nil {
		t.Error(err)
	}

	// It should have created the file in the payload directory.
	_, err := os.Stat(filepath.Join(bag.Path(), "data", expFile))
	if err != nil {
		t.Error("Testing if payload file created:", err)
	}

	// It should have calulated the fixity and put it in the manifest.
	mf := bag.Manifest
	expKey := filepath.Join("data", expFile)
	fx, ok := mf.Data[expKey]
	if !ok {
		t.Error("Unable to find entry in manfest: ", expKey)
	}
	if len(fx) != 32 {
		t.Errorf("Expected %d character fixity but returned: %d", 32, len(fx))
	}
}

func TestAddDir(t *testing.T) {

	// Setup source files to test
	srcDir, _ := ioutil.TempDir("", "_GOTEST_PAYLOAD_SRC_")
	for i := 0; i < 50; i++ {
		fi, _ := ioutil.TempFile(srcDir, "TEST_GO_ADDFILE_")
		fi.WriteString(FIXSTRING)
		fi.Close()
	}
	defer os.RemoveAll(srcDir)

	// Setup the test bag
	bag, _ := setupTestBag("_GOTEST_BAG_ADDDIR_")
	defer os.RemoveAll(bag.Path())

	// It should produce no errors
	if errs := bag.AddDir(srcDir); len(errs) != 0 {
		t.Error(errs)
	}

	// It should produce 50 manifest entries
	if len(bag.Manifest.Data) != 50 {
		t.Error("Expected 50 manifest entries but returned", len(bag.Manifest.Data))
	}
	// It should contain the proper checksums for each file.
	errs := bag.Manifest.RunChecksums()
	for _, err := range errs {
		t.Errorf("%s", err)
	}
}

func TestManifest(t *testing.T) {

	// Setup the test bag
	bag, _ := setupTestBag("_GOTEST_BAG_MANIFEST_")
	defer os.RemoveAll(bag.Path())

	// It should have the expected name and return no error.
	mf := bag.Manifest
	exp := "manifest-md5.txt"
	if filepath.Base(mf.Name()) != exp {
		t.Error("Expected manifest name", exp, "but returned", filepath.Base(mf.Name()))
	}
}

func TestAddTagFile(t *testing.T) {

	// Setup the test bag
	bag, _ := setupTestBag("_GOTEST_BAG_ADDTAGFILE_")
	defer os.RemoveAll(bag.Path())

	// It should throw an error when a bag tagfilename is passed.
	badTagName := "customtag/directory/tag"
	if err := bag.AddTagfile(badTagName); err == nil {
		t.Error("Did not generate an error when trying to add bag tagname:", badTagName)
	}

	// It should not throw an error.
	newTagName := "customtag/directory/tag.txt"
	if err := bag.AddTagfile(newTagName); err != nil {
		t.Error(err)
	}

	// It should be able to lookup the tagfile by name.
	if _, err := bag.TagFile(newTagName); err != nil {
		t.Error(err)
	}

	// Even tagfiles passed as root should be put under the bag.
	oddTagName := "/lookslikeroot/directory/tag.txt"
	if err := bag.AddTagfile(oddTagName); err != nil {
		t.Error(err)
	}

	// It should be able to lookup the tagfile by name.
	if _, err := bag.TagFile(oddTagName); err != nil {
		t.Error(err)
	}

	// It should find the file inside the bag.
	if _, err := os.Stat(filepath.Join(bag.Path(), newTagName)); err != nil {
		t.Error(err)
	}
}

func TestTagFile(t *testing.T) {

	// Setup the test bag
	bag, _ := setupTestBag("_GOTEST_BAG_TAGFILE_")
	defer os.RemoveAll(bag.Path())

	// It should find the tag file by name
	testTagName := "new/tag.txt"
	bag.AddTagfile(testTagName)
	if _, err := bag.TagFile(testTagName); err != nil {
		t.Error(err)
	}

	// It should return an error if asking for a bad tag name.
	badTagName := "/new/tag.txt"
	if _, err := bag.TagFile(badTagName); err == nil {
		t.Error("Bag.TagFile returned results for", badTagName, "when it should not exist.")
	}
}

func TestPath(t *testing.T) {
	// Stup the test bag
	bagName := "_GOTEST_BAG_PATH_"
	bag, _ := setupTestBag(bagName)
	defer os.RemoveAll(bag.Path())

	expPath := filepath.Join(os.TempDir(), bagName)
	if bag.Path() != expPath {
		t.Error("Excpected", bag.Path(), "and", expPath, "to be equal!")
	}
}

func TestSave(t *testing.T) {
	// Setup test bag
	bag, _ := setupTestBag("_GOTEST_BAG_CLOSE_")
	defer os.RemoveAll(bag.Path())

	// Add some data to the manifest and make sure it writes it on close.
	bag.Manifest.Data["data/fakefile.txt"] = "da909ba395016f2a64b04d706520db6afa74fc95"

	// It should not throw an error.
	if errs := bag.Save(); len(errs) != 0 {
		for idx := range errs {
			t.Error(errs[idx])
		}
	}

	// The manifest file should contain data.
	content, err := ioutil.ReadFile(bag.Manifest.Name())
	if err != nil {
		t.Error(err)
	}
	exp := 59 // Length of values entered above and newline.
	if len(content) != 59 {
		t.Error("Expected ", exp, "but found", len(content), "characters written")
	}

	// Add some tagfile data to make sure it writes it on close.
	tfName := "extratagfile.txt"
	bag.AddTagfile("extratagfile.txt")
	tf, _ := bag.TagFile(tfName)
	tf.Data.AddField(*bagins.NewTagField("MyNewField", "This is testdata."))

	// it should not throw an error.
	if errs := bag.Save(); len(errs) != 0 {
		for idx := range errs {
			t.Error(errs[idx])
		}
	}

	// The TagFile should contain data.
	content, err = ioutil.ReadFile(tf.Name())
	if err != nil {
		t.Error(err)
	}
	exp = 10 // Some length the string needs to be abovel
	if len(content) < exp {
		t.Error("Didn't find data in tagfile", tfName, "as expected!")
	}
}

func TestListFiles(t *testing.T) {

	// Setup the test bag.
	bag, _ := setupTestBag("_GOTEST_BAG_LISTFILES_")
	defer os.RemoveAll(bag.Path())

	// Setup the test file to add for the test.
	fi, _ := os.Create((filepath.Join(bag.Path(), "data", "TEST_GO_DATAFILE.txt")))
	fi.WriteString("Test the checksum")
	fi.Close()

	expFiles := make(map[string]bool)
	expFiles["manifest-md5.txt"] = true
	expFiles["bagit.txt"] = true
	expFiles[filepath.Join("data", "TEST_GO_DATAFILE.txt")] = true

	cn, _ := bag.ListFiles()

	for _, fName := range cn {
		if _, ok := expFiles[fName]; !ok {
			t.Error("Unexpected file:", fName)
		}
	}
}
