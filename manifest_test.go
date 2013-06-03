// manifest_test
package bagins_test

import (
	"bagins"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
	"testing"
)

func TestNewManifest(t *testing.T) {
	name := path.Join(os.TempDir(), "_GOTEST_manifest-sha1.txt")
	_, err := bagins.NewManifest(name)
	if err != nil {
		t.Error("Manifest could not be created!", err)
	}
	name = path.Join(os.TempDir(), "_GOTEST_manifest-sha1")
	if err != nil {
		t.Error("NewManifest incorrectly accepting an improperly formatted filename.")
	}
}

func TestAlgoName(t *testing.T) {
	tst := make([]string, 0)
	tst = append(tst, path.Join(os.TempDir(), "_GOTEST_manifest-sha1.txt"))
	tst = append(tst, path.Join(os.TempDir(), "_GOTEST_manifest-md5-sha1.txt"))
	for i := range tst {
		m, _ := bagins.NewManifest(tst[i])
		name, _ := m.AlgoName()
		if name != "sha1" {
			t.Error("AlgoName returned", name, "but expected sha1")
		}
	}
}

func TestRunChecksums(t *testing.T) {
	testFile, _ := ioutil.TempFile("", "_GOTEST_")
	testFile.WriteString("Test the checksum")
	testFile.Close()

	mfst, _ := bagins.NewManifest(path.Join(os.TempDir(), "_GOTEST_manifest-sha1.txt"))
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

func TestManifestCreate(t *testing.T) {
	m, _ := bagins.NewManifest(path.Join(os.TempDir(), "_GOTEST_manifest-sha1.txt"))

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

//func TestManifestName(t *testing.T) {
//	m := new(bagins.Manifest)
//	if name := m.Name(); name != "" {
//		t.Error("Expected empty string for unset BaseName and Algo but returned", m.Name())
//	}

//	// Set only BaseName should still be blank.
//	m.BaseName = "_GOTEST_manifest"
//	if name := m.Name(); name != "" {
//		t.Error("Expected empty string for unset Algo but returned", m.Name())
//	}

//	// Set only Algo should still be blank.
//	m.BaseName, m.Algo = "", "sha1"
//	if name := m.Name(); name != "" {
//		t.Error("Expected empty string for unset BaseName but returned", m.Name())
//	}

//	m.BaseName, m.Algo = "_GOTEST_manifest", "sha1"
//	expected := "_GOTEST_manifest-sha1.txt"
//	if name := m.Name(); name != expected {
//		t.Error("Expected name", expected, "but returned", m.Name())
//	}
//}
