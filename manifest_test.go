// manifest_test
package bagins_test

import (
	"bagins"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
)

func TestRunChecksums(t *testing.T) {
	testFile, _ := ioutil.TempFile("", "_GOTEST_")
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

func TestManifestCreate(t *testing.T) {
	m := bagins.Manifest{Algo: "md5", Data: make(map[string]string), BaseName: "_GOTEST_manifest", Path: os.TempDir()}

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
	m := new(bagins.Manifest)
	if name := m.Name(); name != "" {
		t.Error("Expected empty string for unset BaseName and Algo but returned", m.Name())
	}

	// Set only BaseName should still be blank.
	m.BaseName = "_GOTEST_manifest"
	if name := m.Name(); name != "" {
		t.Error("Expected empty string for unset Algo but returned", m.Name())
	}

	// Set only Algo should still be blank.
	m.BaseName, m.Algo = "", "sha1"
	if name := m.Name(); name != "" {
		t.Error("Expected empty string for unset BaseName but returned", m.Name())
	}

	m.BaseName, m.Algo = "_GOTEST_manifest", "sha1"
	expected := "_GOTEST_manifest-sha1.txt"
	if name := m.Name(); name != expected {
		t.Error("Expected name", expected, "but returned", m.Name())
	}
}
