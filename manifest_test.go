// manifest_test
package bagins_test

import (
	"fmt"
	"github.com/APTrust/bagins"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var test_list map[string]string = map[string]string{
	"md5":    "9e107d9d372bb6826bd81d3542a419d6",
	"sha1":   "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12",
	"sha256": "d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e592",
	"sha512": "07e547d9586f6a73f73fbac0435ed76951218fb7d0c8d788a309d785436bbb642e93a252a954f23912547d1e8a3b5ed6e1bfd7097821233fa0538f3db854fee6",
	"sha224": "730e109bd7a8a32b1cb9d9a09aa2325d2430587ddbc0c38bad911525",
	"sha384": "ca737f1014a48f4c0b6dd43cb177b0afd9e5169367544c494011e3317dbf9a509cb1e5dc1e85a941bbee3d7f2afbc9b1",
}
var test_string string = "The quick brown fox jumps over the lazy dog"

func TestNewManifest(t *testing.T) {
	pth, _ := ioutil.TempDir("", "_GOTEST_MANIFEST")
	defer os.RemoveAll(pth)

	_, err := bagins.NewManifest(pth, "sha1")
	if err != nil {
		t.Error("Manifest could not be created!", err)
	}
}

func TestReadManifest(t *testing.T) {

	// Setup a bad manifest name
	badpth := filepath.Join(os.TempDir(), "__GOTEST__BADMANIFEST_manifest-sha156.txt")
	badfile, err := os.Create(badpth)
	if err != nil {
		t.Error(err)
	}
	badfile.Close()
	defer os.Remove(badfile.Name())

	// It should
	_, errs := bagins.ReadManifest(badpth)
	if len(errs) != 1 {
		t.Error("Did not raise error as expected when trying to read bad manifest filename", badpth)
	}

	// Setup a good manfiest file for tests that should pass.
	exp := make(map[string]string)
	for i := 0; i < 40; i++ {
		check := fmt.Sprintf("%x", rand.Int31())
		fname := fmt.Sprintf("data/testfilename with spaces %d.txt", i)
		exp[fname] = check
	}

	// Setup a good test manifest
	mf, err := bagins.NewManifest(os.TempDir(), "md5")
	if err != nil {
		t.Error(err)
	}
	mf.Data = exp
	err = mf.Create()
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(mf.Name())

	// It should open it and read the values inside without errors.
	m, errs := bagins.ReadManifest(mf.Name())
	if len(errs) != 0 {
		t.Error(errs)
	}
	for fname, check := range exp {
		actual, ok := m.Data[fname]
		if !ok {
			t.Errorf("Expected key %s not found in manifest data", fname)
		}
		if actual != check {
			t.Error("Failed to find file", fname, "in manifest.")
		}
	}
}

func TestRunChecksums(t *testing.T) {
	testFile, _ := ioutil.TempFile("", "_GOTEST_RUNCHECKSUMS.txt")
	testFile.WriteString(test_string)
	testFile.Close()

	mfst, _ := bagins.NewManifest(os.TempDir(), "sha1")
	mfst.Data[filepath.Base(testFile.Name())] = test_list["sha1"]
	errList := mfst.RunChecksums()

	// Checksum for file should now be generated.
	for _, err := range errList {
		t.Error(err)
	}

	// Check that it throws an error if mismatch checksum.
	mfst.Data[testFile.Name()] = "frodo lives!"
	errList = mfst.RunChecksums()
	if len(errList) == 0 {
		t.Error("Invalid Checksums not being detected!")
	}
	os.Remove(testFile.Name()) // Remove the test file.
}

func TestManifestCreate(t *testing.T) {
	m, _ := bagins.NewManifest(os.TempDir(), "sha1")

	testFiles := make([]*os.File, 3)
	for idx := range testFiles {
		testFiles[idx], _ = ioutil.TempFile("", "_GOTEST_")
		testFiles[idx].WriteString(strings.Repeat("test ", rand.Intn(50)))
		m.Data[testFiles[idx].Name()] = ""
		testFiles[idx].Close()
	}

	m.RunChecksums()
	m.Create()

	// Clean it all up.
	for idx := range testFiles {
		os.Remove(testFiles[idx].Name())
	}
	os.Remove(m.Name())
}

func TestManifestName(t *testing.T) {

	// Set only Algo should still be blank.
	m, err := bagins.NewManifest(os.TempDir(), "SHA1")
	if err != nil {
		t.Error(err)
	}
	exp := filepath.Join(os.TempDir(), "manifest-sha1.txt")
	if name := m.Name(); name != exp {
		t.Error("Expected mainfest name %s but returned %s", exp, m.Name())
	}
}

func TestManifestToString(t *testing.T) {
	m, _ := bagins.NewManifest(os.TempDir(), "sha1")
	m.Data["FileOne.txt"] = fmt.Sprintf("CHECKSUM 0001")
	m.Data["FileTwo.txt"] = fmt.Sprintf("CHECKSUM 0002")
	m.Data["FileThree.txt"] = fmt.Sprintf("CHECKSUM 0003")

	output := m.ToString()
	lines := []string {
		"CHECKSUM 0001 FileOne.txt\n",
		"CHECKSUM 0002 FileTwo.txt\n",
		"CHECKSUM 0003 FileThree.txt\n",
	}

	for _, line := range lines {
		if !strings.Contains(output, line) {
			t.Errorf("Manifest.ToString() did not return line %s", line)
		}
	}
	expectedLength := len(lines[0]) + len(lines[1]) + len(lines[2])
	if len(output) != expectedLength {
		t.Errorf("Manifest.ToString() returned %d characters, expected %d",
			len(output), expectedLength)
	}
}
