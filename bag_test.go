package bagins_test

import (
	//	"fmt"
	"github.com/APTrust/bagins"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	FIXSTRING string = "The quick brown fox jumps over the lazy dog."
	FIXVALUE  string = "e4d909c290d0fb1ca068ffaddf22cbd0" // md5 of string above
)

func setupTestBag(bagName string) (*bagins.Bag, error) {
	bag, err := bagins.NewBag(os.TempDir(), bagName, []string{"md5"}, false)
	if err != nil {
		return nil, err
	}
	return bag, nil
}

// Setups up a bag with some custom tag files.
func setupCustomBag(bagName string) (*bagins.Bag, error) {
	bag, err := bagins.NewBag(os.TempDir(), bagName, []string{"md5", "sha256"}, true)
	if err != nil {
		return nil, err
	}
	bag.AddTagfile("bag-info.txt")
	bagInfo, _ := bag.TagFile("bag-info.txt")
	bagInfo.Data.SetFields([]bagins.TagField{
		*bagins.NewTagField("Source-Organization", "APTrust"),
		*bagins.NewTagField("Bagging-Date", "2016-06-01"),
		*bagins.NewTagField("Bag-Count", "1"),
		*bagins.NewTagField("Internal-Sender-Description", "This is a test bag with no content."),
		*bagins.NewTagField("Internal-Sender-Identification", "Bag XYZ"),
	})
	bag.AddTagfile("aptrust-info.txt")
	aptrustInfo, _ := bag.TagFile("aptrust-info.txt")
	aptrustInfo.Data.SetFields([]bagins.TagField{
		*bagins.NewTagField("Title", "APTrust Generic Test Bag"),
		*bagins.NewTagField("Rights", "Consortia"),
	})
	errors := bag.Save()
	if errors != nil && len(errors) > 0 {
		return nil, errors[0]
	}
	return bag, nil
}

// Setups up a bag with custom tag files, a custom tag directory,
// and a tag manifest.
func setupTagfileBag(bagName string) (*bagins.Bag, error) {
	bag, err := setupCustomBag(bagName)
	if err != nil {
		return nil, err
	}

	// Tag file in top-level directory
	bag.AddTagfile("laser-tag.txt")
	customTagFile1, _ := bag.TagFile("laser-tag.txt")
	customTagFile1.Data.SetFields([]bagins.TagField{
		*bagins.NewTagField("tic", "tac"),
		*bagins.NewTagField("tick", "tock"),
	})

	// Tag files in custom directory
	bag.AddTagfile("custom-tags/player-stats.txt")
	customTagFile2, _ := bag.TagFile("custom-tags/player-stats.txt")
	customTagFile2.Data.SetFields([]bagins.TagField{
		*bagins.NewTagField("Batting-Average", ".340"),
		*bagins.NewTagField("What-Time-Is-It", "2016-06-01T12:00:00Z"),
		*bagins.NewTagField("ERA", "1.63"),
		*bagins.NewTagField("Bats", "Left"),
		*bagins.NewTagField("Throws", "Right"),
	})
	bag.AddTagfile("custom-tags/tv-schedule.txt")
	customTagFile3, _ := bag.TagFile("custom-tags/tv-schedule.txt")
	customTagFile3.Data.SetFields([]bagins.TagField{
		*bagins.NewTagField("3:00PM", "House Party"),
		*bagins.NewTagField("4:00PM", "Dexter"),
	})

	errors := bag.Save()
	if errors != nil && len(errors) > 0 {
		return nil, errors[0]
	}
	return bag, nil
}

func TestNewBag(t *testing.T) {

	// It should raise an error if the destination dir does not exist.
	badLocation := filepath.Join(os.TempDir(), "/GOTESTNOT_EXISTs/")
	_, err := bagins.NewBag(badLocation, "_GOFAILBAG_", []string{"md5"}, false)
	if err == nil {
		t.Error("NewBag function does not recognize when a directory does not exist!")
	}

	// It should raise an error if the bag already exists.
	os.MkdirAll(filepath.Join(badLocation, "_GOFAILBAG_"), 0766)
	defer os.RemoveAll(badLocation)

	_, err = bagins.NewBag(badLocation, "_GOFAILBAG_", []string{"md5"}, false)
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
	if _, err := bagins.ReadBag(badPath, []string{}); err == nil {
		t.Errorf("Path %s not detected as bad as expected.", badPath)
	}

	// It should return an error if it isn't passed a path to a directory.
	fi, _ := ioutil.TempFile("", "TEST_GO_READBAG_")
	fi.WriteString("Test file please delete.")
	fi.Close()
	defer os.Remove(fi.Name())

	if _, err := bagins.ReadBag(fi.Name(), []string{}); err == nil {
		t.Errorf("Readbag should thrown an error when trying to open a file: %s", fi.Name())
	}

	// It should return an error if the directory does not contain a data subdirectory.
	pDir, _ := ioutil.TempDir("", "_GOTEST_ReadBag_Payload_")
	defer os.RemoveAll(pDir)

	if _, err := bagins.ReadBag(pDir, []string{}); err == nil {
		t.Errorf("Not returning expected error when directory has no data subdirectory for %s", pDir)
	}

	os.Mkdir(filepath.Join(pDir, "data"), os.ModePerm) // Set up data directory for later tests.

	// It should return an error if there is no manifest file.
	if _, err := bagins.ReadBag(pDir, []string{}); err == nil {
		t.Errorf("Not returning expected error when no manifest file is present in %s", pDir)
	}

	// It should return an error if it has a bad manifest name.
	ioutil.WriteFile(filepath.Join(pDir, "manifest-sha404.txt"), []byte{}, os.ModePerm)
	if _, err := bagins.ReadBag(pDir, []string{}); err == nil {
		t.Errorf("Not returning expected error when a bad manifest filename is only option %s", pDir)
	}
	os.Remove(filepath.Join(pDir, "manifest-sha404.txt"))

	// It should return a bag if a valid manifest and data directory exist.
	ioutil.WriteFile(filepath.Join(pDir, "manifest-sha256.txt"), []byte{}, os.ModePerm)
	if _, err := bagins.ReadBag(pDir, []string{}); err != nil {
		t.Errorf("Unexpected error when trying to read raw bag with valid data and manifest: %s", err)
	}

	// It should read and return a valid bag object reading with a baginfo.txt tagfile.
	bagName := "__GO_READBAG_TEST__"
	bagPath := filepath.Join(os.TempDir(), bagName)
	tb, err := setupTestBag(bagName)
	defer os.RemoveAll(bagPath)
	if err != nil {
		t.Errorf("%s", err)
	}
	tb.Save()

	testBag, err := bagins.ReadBag(bagPath, []string{"bagit.txt"})
	if err != nil {
		t.Errorf("Unexpected error reading test bag: %s", err)
	}
	baginfo, err := testBag.TagFile("bagit.txt")
	if err != nil {
		t.Errorf("Unexpected error reading bagit.txt file: %s", err)
	}
	if baginfo == nil {
		t.Errorf("Baginfo unexpectedly nil.")
	}
}

func TestReadCustomBag(t *testing.T) {
	// Setup File to test.
	fi, _ := ioutil.TempFile("", "TEST_READ_CUSTOM_BAG_FILE.txt")
	fi.WriteString(FIXSTRING)
	fi.Close()
	defer os.Remove(fi.Name())

	// Setup Custom Bag
	bagName := "__GO_TEST_READ_CUSTOM_BAG__"
	bagPath := filepath.Join(os.TempDir(), bagName)
	defer os.RemoveAll(bagPath)
	bag, err := setupCustomBag(bagName)
	if err != nil {
		t.Errorf("Unexpected error setting up custom bag: %s", err)
	}
	bag.AddFile(fi.Name(), fi.Name())
	bag.Save()
	defer os.RemoveAll(bagPath)

	rBag, err := bagins.ReadBag(bag.Path(), []string{"bag-info.txt", "aptrust-info.txt"})
	if err != nil {
		t.Errorf("Unexpected error reading custom bag: %s", err)
	}

	bagInfo, err := rBag.TagFile("bag-info.txt")
	if err != nil {
		t.Errorf("Error finding bag-info.txt tag file: %s", err)
	}
	if len(bagInfo.Data.Fields()) != 5 {
		t.Errorf("Expected 5 fields in bag-info.txt but returned %d", len(bagInfo.Data.Fields()))
	}

	aptrustInfo, err := rBag.TagFile("aptrust-info.txt")
	if err != nil {
		t.Errorf("Error finding aptrust-info.txt tag file: %s", err)
	}
	if len(aptrustInfo.Data.Fields()) != 2 {
		t.Errorf("Expected 2 fields in aptrust-info.txt but returned %d", len(aptrustInfo.Data.Fields()))
	}

	// Check payload manifests
	payloadManifests := rBag.GetManifests(bagins.PayloadManifest)
	if len(payloadManifests) != 2 {
		t.Errorf("Expected 2 payload manifests, got %d", len(payloadManifests))
	} else {
		if payloadManifests[0].Algorithm() != "md5" {
			t.Errorf("Expected first manifest to be md5, got %s", payloadManifests[0].Algorithm())
		}
		// Payload md5 manifest should have one entry
		if len(payloadManifests[0].Data) != 1 {
			t.Errorf("Payload manifest should have one entry, found %s", len(payloadManifests[0].Data))
		}
		for key, value := range payloadManifests[0].Data {
			dataFilePath := filepath.Join("data", fi.Name())
			if key != dataFilePath {
				t.Errorf("Missing expected manifest entry for %s. Got %s", dataFilePath, key)
			}
			if value != "e4d909c290d0fb1ca068ffaddf22cbd0" {
				t.Errorf("Incorrect md5 checksum. Got %s", value)
			}
		}

		if payloadManifests[1].Algorithm() != "sha256" {
			t.Errorf("Expected first manifest to be sha256, got %s", payloadManifests[1].Algorithm())
		}
		// Payload sha256 manifest should have one entry
		if len(payloadManifests[1].Data) != 1 {
			t.Errorf("Payload manifest should have one entry, found %s", len(payloadManifests[1].Data))
		}
		for key, value := range payloadManifests[1].Data {
			dataFilePath := filepath.Join("data", fi.Name())
			if key != dataFilePath {
				t.Errorf("Missing expected manifest entry for %s. Got %s", dataFilePath, key)
			}
			if value != "ef537f25c895bfa782526529a9b63d97aa631564d5d789c2b765448c8635fb6c" {
				t.Errorf("Incorrect md5 checksum. Got %s", value)
			}
		}
	}
}

func TestReadTagFileBag(t *testing.T) {
	// Setup File to test.
	testFileName := "TEST_READ_TAGFILE_BAG_FILE.txt"
	fi, _ := ioutil.TempFile("", testFileName)
	fi.WriteString(FIXSTRING)
	fi.Close()
	defer os.Remove(fi.Name())

	// Setup bag with custom tag files
	bagName := "__GO_TEST_READ_TAGFILE_BAG__"
	bagPath := filepath.Join(os.TempDir(), bagName)
	defer os.RemoveAll(bagPath)

	bag, err := setupTagfileBag(bagName)
	if err != nil {
		t.Errorf("Unexpected error setting up tagfile bag: %s", err)
	}
	bag.AddFile(fi.Name(), testFileName)
	bag.Save()

	rBag, err := bagins.ReadBag(bag.Path(), []string{"bagit.txt", "bag-info.txt", "aptrust-info.txt"})
	if err != nil {
		t.Errorf("Unexpected error reading custom bag: %s", err)
	}

	bagInfo, err := rBag.TagFile("bag-info.txt")
	if err != nil {
		t.Errorf("Error finding bag-info.txt tag file: %s", err)
	}
	if len(bagInfo.Data.Fields()) != 5 {
		t.Errorf("Expected 5 fields in bag-info.txt but returned %d", len(bagInfo.Data.Fields()))
	}

	aptrustInfo, err := rBag.TagFile("aptrust-info.txt")
	if err != nil {
		t.Errorf("Error finding aptrust-info.txt tag file: %s", err)
	}
	if len(aptrustInfo.Data.Fields()) != 2 {
		t.Errorf("Expected 2 fields in aptrust-info.txt but returned %d", len(aptrustInfo.Data.Fields()))
	}

	// Check payload manifests
	payloadManifests := rBag.GetManifests(bagins.PayloadManifest)
	if len(payloadManifests) != 2 {
		t.Errorf("Expected 2 payload manifests, got %d", len(payloadManifests))
	} else {
		if payloadManifests[0].Algorithm() != "md5" {
			t.Errorf("Expected first manifest to be md5, got %s", payloadManifests[0].Algorithm())
		}
		// Payload md5 manifest should have one entry
		if len(payloadManifests[0].Data) != 1 {
			t.Errorf("Payload manifest should have one entry, found %s", len(payloadManifests[0].Data))
		}
		for key, value := range payloadManifests[0].Data {
			dataFilePath := filepath.Join("data", testFileName)
			if key != dataFilePath {
				t.Errorf("Missing expected manifest entry for %s. Got %s", dataFilePath, key)
			}
			if value != "e4d909c290d0fb1ca068ffaddf22cbd0" {
				t.Errorf("Incorrect md5 checksum. Got %s", value)
			}
		}

		if payloadManifests[1].Algorithm() != "sha256" {
			t.Errorf("Expected first manifest to be sha256, got %s", payloadManifests[1].Algorithm())
		}
		// Payload sha256 manifest should have one entry
		if len(payloadManifests[1].Data) != 1 {
			t.Errorf("Payload manifest should have one entry, found %s", len(payloadManifests[1].Data))
		}
		for key, value := range payloadManifests[1].Data {
			dataFilePath := filepath.Join("data", testFileName)
			if key != dataFilePath {
				t.Errorf("Missing expected manifest entry for %s. Got %s", dataFilePath, key)
			}
			if value != "ef537f25c895bfa782526529a9b63d97aa631564d5d789c2b765448c8635fb6c" {
				t.Errorf("Incorrect md5 checksum. Got %s", value)
			}
		}
	}

	// Check that tag files exist on disk
	aptrustInfoFile := "aptrust-info.txt"
	bagInfoFile := "bag-info.txt"
	bagItFile := "bagit.txt"
	laserTagFile := "laser-tag.txt"
	playerStatsFile := filepath.Join("custom-tags", "player-stats.txt")
	tvScheduleFile := filepath.Join("custom-tags", "tv-schedule.txt")
	manifestMd5 := "manifest-md5.txt"
	manifestSha256 := "manifest-sha256.txt"

	tagFiles := []string{aptrustInfoFile, bagInfoFile, bagItFile,
		laserTagFile, playerStatsFile, tvScheduleFile}
	for _, tf := range tagFiles {
		absPath := filepath.Join(bagPath, tf)
		_, err = os.Stat(absPath)
		if err != nil && os.IsNotExist(err) {
			t.Errorf("Tag file is not written to disk at %s", tf)
		}
	}

	// Make sure the bag knows they're there too
	files, err := rBag.ListFiles()
	if err != nil {
		t.Errorf("ListFiles() returned error %v", err)
	}
	expectedFiles := []string{
		"aptrust-info.txt",
		"bag-info.txt",
		"bagit.txt",
		"custom-tags/player-stats.txt",
		"custom-tags/tv-schedule.txt",
		"data/TEST_READ_TAGFILE_BAG_FILE.txt",
		"laser-tag.txt",
		"manifest-md5.txt",
		"manifest-sha256.txt",
		"tagmanifest-md5.txt",
		"tagmanifest-sha256.txt",
	}
	for _, expectedFile := range expectedFiles {
		if !sliceContains(files, expectedFile) {
			t.Errorf("ListFiles did not return file %s", expectedFile)
		}
	}

	// Make sure the bag knows about the parsed tag files
	expectedParsedFiles := []string{"bagit.txt", "bag-info.txt", "aptrust-info.txt"}
	parsedTagFiles := rBag.ListTagFiles()
	if len(parsedTagFiles) != len(expectedParsedFiles) {
		t.Errorf("Expected %d parsed tag files, got %d",
			len(parsedTagFiles), len(expectedParsedFiles))
	}

	for _, expected := range expectedParsedFiles {
		if !sliceContains(parsedTagFiles, expected) {
			t.Errorf("ListTagFiles() did not return file %s", expected)
		}
	}

	// Make sure the bag knows about unparsed tag files
	expectedUnparsed := []string{
		"custom-tags/player-stats.txt",
		"custom-tags/tv-schedule.txt",
		"laser-tag.txt"}
	unparsedTagFiles, err := rBag.UnparsedTagFiles()
	if err != nil {
		t.Errorf("UnparsedTagFiles() returned unexpected error: %v", err)
	}

	if len(unparsedTagFiles) != len(expectedUnparsed) {
		t.Errorf("Expected %d parsed tag files, got %d",
			len(unparsedTagFiles), len(expectedUnparsed))
	}

	for _, expected := range expectedUnparsed {
		if !sliceContains(unparsedTagFiles, expected) {
			t.Errorf("UnparsedTagFiles() did not return file %s", expected)
		}
	}

	// Check tag manifests
	tagManifests := rBag.GetManifests(bagins.TagManifest)
	if len(tagManifests) != 2 {
		t.Errorf("Expected 2 tag manifests, got %d", len(tagManifests))
	} else {
		if tagManifests[0].Algorithm() != "md5" {
			t.Errorf("Expected first manifest to be md5, got %s", tagManifests[0].Algorithm())
		}
		// Tag md5 manifest should have six entries
		if len(tagManifests[0].Data) != 8 {
			t.Errorf("Tag manifest should have 8 entries, found %d", len(tagManifests[0].Data))
		}

		// Check the fixity values
		md5Entries := make(map[string]string, 6)
		md5Entries[aptrustInfoFile] = "6dd711392d4661322acc469a30565f68"
		md5Entries[bagInfoFile] = "88190858fd93609ae51ca1f06ee575f1"
		md5Entries[bagItFile] = "ada799b7e0f1b7a1dc86d4e99df4b1f4"
		md5Entries[laserTagFile] = "29251712228b36927c43157fe5808552"
		md5Entries[playerStatsFile] = "dfa872f6da2af8087bea5f7ab1dbc1fa"
		md5Entries[tvScheduleFile] = "118df3be000eae34d6e6dbf7f56c649b"
		md5Entries[manifestMd5] = "4c5a8c217cf51fb419e425ddc2f433ee"
		md5Entries[manifestSha256] = "edbc9b8dabe4c894d22cc42d5268867b"

		for key, expectedValue := range md5Entries {
			actualValue := tagManifests[0].Data[key]
			if actualValue != expectedValue {
				t.Errorf("For tag file %s, expected md5 %s, but got %s",
					key, expectedValue, actualValue)
			}
		}

		if tagManifests[1].Algorithm() != "sha256" {
			t.Errorf("Expected first manifest to be sha256, got %s", tagManifests[1].Algorithm())
		}
		// Tag sha256 manifest should have 8 entries (2 are for payload manifests)
		if len(tagManifests[1].Data) != 8 {
			t.Errorf("Tag manifest should have 8 entries, found %d", len(tagManifests[1].Data))
		}

		// Check fixity values
		// Will these checksums break on Windows, where end-of-line is CRLF?
		sha256Entries := make(map[string]string, 6)
		sha256Entries[aptrustInfoFile] =
			"ffe2ab04b87db85886fcfd013c9f09e094b636ca233cd0cbbd1ea300e7a5352c"
		sha256Entries[bagInfoFile] =
			"f0ce035c2ee789a7f8821d6f174a75619c575eea0311c47d03149807d252804d"
		sha256Entries[bagItFile] =
			"49b477e8662d591f49fce44ca5fc7bfe76c5a71f69c85c8d91952a538393e5f4"
		sha256Entries[laserTagFile] =
			"163be000df169eafd84fa0cef6028a4711e53bd3abf9e8c54603035bb92bda95"
		sha256Entries[playerStatsFile] =
			"83137fc6d88212250153bd713954da1d1c5a69c57a55ff97cac07ca6db7ec34d"
		sha256Entries[tvScheduleFile] =
			"fbf223502fe7f470363346283620401d04e77fe43a9a74faa682eebe28417e7c"
		sha256Entries[manifestMd5] =
			"9aab27bb0d2d75d7ac2c26908e2ca85e7121f106445318f7def024f4b520bec2"
		sha256Entries[manifestSha256] =
			"2a9e5d86070459a652fdc2ce13f5358ede42b10a3e0580e149b0d3df938ffe30"

		for key, expectedValue := range sha256Entries {
			actualValue := tagManifests[1].Data[key]
			if actualValue != expectedValue {
				t.Errorf("For tag file %s, expected sha256 %s, but got %s",
					key, expectedValue, actualValue)
			}
		}
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
	if bag.Manifests == nil || len(bag.Manifests) == 0 {
		t.Error("Bag manifest is missing")
		return
	}
	mf := bag.Manifests[0]
	expKey := filepath.Join("data", expFile)
	fx, ok := mf.Data[expKey]
	if !ok {
		t.Error("Unable to find entry in manfest: ", expKey)
	}
	if len(fx) != 32 {
		t.Errorf("Expected %d character fixity but returned: %d", 32, len(fx))
	}
}

func TestAddCustomTagfile(t *testing.T) {
	// Setup File to test.
	fi, _ := ioutil.TempFile("", "TEST_READ_CUSTOM_BAG_FILE.txt")
	fi.WriteString(FIXSTRING)
	fi.Close()
	defer os.Remove(fi.Name())

	// Setup Custom Bag
	bagName := "__GO_TEST_ADD_CUSTOM_TAG_FILE__"
	bagPath := filepath.Join(os.TempDir(), bagName)
	defer os.RemoveAll(bagPath)
	bag, err := setupCustomBag(bagName)
	if err != nil {
		t.Errorf("Unexpected error setting up custom bag: %s", err)
	}
	bag.AddCustomTagfile(fi.Name(), "custom-tags/in_manifest.txt", true)
	bag.AddCustomTagfile(fi.Name(), "custom-tags/not_in_manifest.txt", false)
	bag.Save()
	defer os.RemoveAll(bagPath)

	rBag, err := bagins.ReadBag(bag.Path(), []string{"bag-info.txt", "aptrust-info.txt"})
	if err != nil {
		t.Errorf("Unexpected error reading custom bag: %s", err)
	}

	files, err := rBag.ListFiles()
	if err != nil {
		t.Errorf("Error listing bag files: %v", err)
	}

	// Make sure the file exists
	if !sliceContains(files, "custom-tags/in_manifest.txt") {
		t.Errorf("Custom tag file 'custom-tags/in_manifest.txt' is not in the bag")
	}
	if !sliceContains(files, "custom-tags/not_in_manifest.txt") {
		t.Errorf("Custom tag file 'custom-tags/not_in_manifest.txt' is not in the bag")
	}

	// First file should be in the tag manifests. Second file should not.
	tagManifests := rBag.GetManifests(bagins.TagManifest)
	for _, tagManifest := range tagManifests {
		if _, exists := tagManifest.Data["custom-tags/in_manifest.txt"]; exists == false {
			t.Errorf("File 'custom-tags/in_manifest.txt' is missing from tagmanifest-%s.txt",
				tagManifest.Algorithm())
		}
		if _, exists := tagManifest.Data["custom-tags/not_in_manifest.txt"]; exists == true {
			t.Errorf("File 'custom-tags/not_in_manifest.txt' is should not have an entry in "+
				"tagmanifest-%s.txt, but it does", tagManifest.Algorithm())
		}
	}
}

func TestListTagFiles(t *testing.T) {
	// Setup Custom Bag
	bagName := "__GO_TEST_LIST_TAG_FILES__"
	bagPath := filepath.Join(os.TempDir(), bagName)
	defer os.RemoveAll(bagPath)
	bag, err := setupCustomBag(bagName)
	if err != nil {
		t.Errorf("Unexpected error setting up custom bag: %s", err)
	}

	expected := []string{"bagit.txt", "bag-info.txt", "aptrust-info.txt"}
	if len(bag.ListTagFiles()) != len(expected) {
		t.Errorf("Expected %d tag files but returned %d", len(expected), len(bag.ListTagFiles()))
	}
	for _, name := range expected {
		if _, err := bag.TagFile(name); err != nil {
			t.Errorf("Error getting tag file %s: %s", name, err)
		}
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
	bag, err := setupTestBag("_GOTEST_BAG_ADDDIR_")
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer os.RemoveAll(bag.Path())

	// It should produce no errors
	if errs := bag.AddDir(srcDir); len(errs) != 0 {
		t.Error(errs)
	}

	// It should produce 50 manifest entries
	if bag.Manifests == nil || len(bag.Manifests) == 0 {
		t.Error("Bag manifest is missing")
		return
	}
	manifest := bag.Manifests[0]
	if len(manifest.Data) != 50 {
		t.Error("Expected 50 manifest entries but returned", len(manifest.Data))
	}
	// It should contain the proper checksums for each file.
	errs := manifest.RunChecksums()
	for _, err := range errs {
		t.Errorf("%s", err)
	}
}

func TestManifest(t *testing.T) {

	// Setup the test bag
	bag, _ := setupTestBag("_GOTEST_BAG_MANIFEST_")
	defer os.RemoveAll(bag.Path())

	// It should have the expected name and return no error.
	if bag.Manifests == nil || len(bag.Manifests) == 0 {
		t.Error("Bag manifest is missing")
		return
	}
	mf := bag.Manifests[0]
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
	bag.Manifests[0].Data["data/fakefile.txt"] = "da909ba395016f2a64b04d706520db6afa74fc95"

	// It should not throw an error.
	if errs := bag.Save(); len(errs) != 0 {
		for idx := range errs {
			t.Error(errs[idx])
		}
	}

	// The manifest file should contain data.
	content, err := ioutil.ReadFile(bag.Manifests[0].Name())
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

func sliceContains(list []string, item string) bool {
	for _, value := range list {
		if value == item {
			return true
		}
	}
	return false
}
