// manifest_test
package bagins_test

import (
	"fmt"
	"github.com/APTrust/bagins"
	"github.com/APTrust/bagins/bagutil"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewManifest(t *testing.T) {
	pth, _ := ioutil.TempDir("", "_GOTEST_MANIFEST")
	defer os.RemoveAll(pth)

	chk, _ := bagutil.NewCheckByName("sha1")
	_, err := bagins.NewManifest(pth, chk)
	if err != nil {
		t.Error("Manifest could not be created!", err)
	}
}

func TestReadManifest(t *testing.T) {
	exp := make(map[string]string)
	for i := 0; i < 40; i++ {
		check := fmt.Sprintf("%x", rand.Int31())
		fname := fmt.Sprintf("data/testfilename_%d.txt", i)
		exp[fname] = check
	}

	// Setup a test manifest
	h, err := bagutil.NewCheckByName("md5")
	if err != nil {
		t.Error(err)
	}
	mf, err := bagins.NewManifest(os.TempDir(), h)
	if err != nil {
		t.Error(err)
	}
	mf.Data = exp
	err = mf.Create()
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(mf.Name())

	// Open it and read the values.
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

func TestGetAlgoName(t *testing.T) {
	tst := make([]string, 0)
	tst = append(tst, filepath.Join(os.TempDir(), "_GOTEST_manifest-sha1.txt"))
	tst = append(tst, filepath.Join(os.TempDir(), "_GOTEST_manifest-md5-sha1.txt"))
	for i := range tst {
		name, _ := bagins.GetAlgoName(tst[i])
		if name != "sha1" {
			t.Error("AlgoName returned", name, "but expected sha1")
		}
	}
}

func TestRunChecksums(t *testing.T) {
	testFile, _ := ioutil.TempFile("", "_GOTEST_")
	testFile.WriteString("Test the checksum")
	testFile.Close()

	chk, _ := bagutil.NewCheckByName("sha1")
	mfst, _ := bagins.NewManifest(os.TempDir(), chk)
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
	chk, _ := bagutil.NewCheckByName("sha1")
	m, _ := bagins.NewManifest(os.TempDir(), chk)

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
	h, _ := bagutil.NewCheckByName("SHA1")
	m, err := bagins.NewManifest(os.TempDir(), h)
	if err != nil {
		t.Error(err)
	}
	exp := filepath.Join(os.TempDir(), "manifest-sha1.txt")
	if name := m.Name(); name != exp {
		t.Error("Expected mainfest name %s but returned %s", exp, m.Name())
	}
}
