package bagins_test

import (
	"crypto/md5"
	"github.com/APTrust/bagins"
	"github.com/APTrust/bagins/bagutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewPayload(t *testing.T) {

	tmpPyld := filepath.Join(os.TempDir(), "_GOTEST_NewPayload_")

	// Check for failure on non-existant directory.
	_, err := bagins.NewPayload(tmpPyld)
	if err == nil {
		t.Errorf("Unexpected error return checking for non-existed directory: %s", err)
	}

	// Check for positive return when directory exists.
	pth, err := ioutil.TempDir("", "_GOTEST_NewPayload_")
	if err != nil {
		t.Errorf("Unexpcted error creating temporary directory: %s", err)
	}
	tstDir, err := os.Stat(pth)
	if err != nil {
		t.Errorf("Reading %s returned an error: %s", pth, err)
	}
	if !tstDir.IsDir() {
		t.Errorf("Payload dir %s is not a valid directory", pth)
	}

	// Clean it up.
	os.RemoveAll(pth)
}

func TestPayloadName(t *testing.T) {
	pDir, _ := ioutil.TempDir("", "_GOTEST_PayloadName_")
	defer os.Remove(pDir)

	p, _ := bagins.NewPayload(pDir)

	if pDir != p.Name() {
		t.Errorf("Payload name %s did not equal expected %s", p.Name(), pDir)
	}
}

func TestPayloadAdd(t *testing.T) {
	pDir, _ := ioutil.TempDir("", "_GOTEST_PayloadAdd_")
	m, _ := bagins.NewManifest(os.TempDir(), "md5")
	defer os.RemoveAll(pDir)

	p, err := bagins.NewPayload(pDir)
	if err != nil {
		t.Error(err)
	}

	testFile, _ := ioutil.TempFile("", "_GO_PayloadAdd_TESTFILE_")
	testFile.WriteString("Test the checksum")
	testFile.Close()
	defer os.Remove(testFile.Name())

	chkSum, err := p.Add(testFile.Name(), filepath.Base(testFile.Name()), []*bagins.Manifest{m})
	if err != nil {
		t.Error(err)
	}
	exp := "92d7a9f0f4a30ca782dcae5fe83ca7eb"
	if exp != chkSum["md5"] {
		t.Error("Checksum", chkSum["md5"], "did not match", exp)
	}
}

// Make sure that when we add a file to the payload
// that is already in the payload directory, it doesn't
// get clobbered.
func TestPayloadAddInPlace(t *testing.T) {
	pDir, _ := ioutil.TempDir("", "_GOTEST_PayloadAdd_")
	m, _ := bagins.NewManifest(os.TempDir(), "md5")
	defer os.RemoveAll(pDir)

	p, err := bagins.NewPayload(pDir)
	if err != nil {
		t.Error(err)
	}

	//testFile, _ := ioutil.TempFile("", "_GO_PayloadAdd_TESTFILE_")
	testFile, err := os.Create(filepath.Join(pDir, "_GO_PayloadAdd_TESTFILE_"))
	if err != nil {
		t.Error(err)
	}
	testFile.WriteString("Test the checksum")
	testFile.Close()
	defer os.Remove(testFile.Name())

	chkSum, err := p.Add(testFile.Name(), filepath.Base(testFile.Name()), []*bagins.Manifest{m})
	if err != nil {
		t.Error(err)
	}
	exp := "92d7a9f0f4a30ca782dcae5fe83ca7eb"
	if exp != chkSum["md5"] {
		t.Error("Checksum", chkSum["md5"], "did not match", exp)
	}
}


func TestPayloadAddAll(t *testing.T) {
	// Setup directories to test on
	srcDir, _ := ioutil.TempDir("", "_GOTEST_PayloadAddAll_SRCDIR_")
	defer os.RemoveAll(srcDir)
	pDir, _ := ioutil.TempDir("", "_GOTEST_PayloadAddAll_")
	defer os.RemoveAll(pDir)

	m, _ := bagins.NewManifest(os.TempDir(), "md5")

	// Setup test files
	for i := 0; i < 100; i++ {
		tstFile, _ := ioutil.TempFile(srcDir, "_GOTEST_PayloadAddAll_FILE_")
		tstFile.WriteString("Test the checksum")
		tstFile.Close()
	}

	p, _ := bagins.NewPayload(pDir)
	checksums, errs := p.AddAll(srcDir, []*bagins.Manifest{m})

	// It should not return an error.
	if errs != nil {
		t.Errorf("Add all returned %d errors", len(errs))
	}
	// It should have fixity values for 100 files
	if len(checksums) != 100 {
		t.Errorf("Expected 100 fixity values but returned %d", len(checksums))
	}

	for key := range checksums {
		fileChk, err := bagutil.FileChecksum(filepath.Join(p.Name(), key), md5.New())
		if err != nil {
			t.Errorf(" %s", err)
		}
		if checksums[key]["md5"] != fileChk {
			t.Error("Expected", checksums[key]["md5"], "but returned", fileChk)
		}
	}

}

func TestPayloadOctetStreamSum(t *testing.T) {
	// Setup Test directory
	pDir, _ := ioutil.TempDir("", "_GOTEST_PayloadOctetStreamSum_")
	defer os.RemoveAll(pDir)

	// Setup test files
	for i := 0; i < 100; i++ {
		tstFile, _ := ioutil.TempFile(pDir, "_GOTEST_PayloadOctetStreamSum_FILE_")
		tstFile.WriteString("Test the checksum")
		tstFile.Close()
	}

	p, _ := bagins.NewPayload(pDir)
	sum, count := p.OctetStreamSum()

	if sum != 1700 {
		t.Error("Sum of octets expected to be 1700 but returned", sum)
	}
	if count != 100 {
		t.Error("Count of files expected to be 100 but returned", count)
	}
}

func BenchmarkPayload(b *testing.B) {
	srcDir, _ := ioutil.TempDir("", "_GOTEST_BenchmarkPayload_SRCDIR_")
	defer os.RemoveAll(srcDir)
	pDir, _ := ioutil.TempDir("", "_GOTEST_BenchmarkPayload_Payload_")
	defer os.RemoveAll(pDir)

	m, _ := bagins.NewManifest(os.TempDir(), "md5")

	// Make src temp test files
	for i := 0; i < 300; i++ {
		tstFile, _ := ioutil.TempFile(srcDir, "_GOTEST_BenchmarkPayload_FILE_")
		tstFile.WriteString(strings.Repeat("Test the checksum. ", 500000)) // produces ~9 meg text file.
		tstFile.Close()
	}

	b.ResetTimer()

	p, _ := bagins.NewPayload(pDir)

	checksums, err := p.AddAll(srcDir, []*bagins.Manifest{m})
	if err != nil {
		b.Error(err)
	}

	b.StopTimer()

	// Make sure the actual values check out.
	for key := range checksums {
		fileChk, err := bagutil.FileChecksum(filepath.Join(p.Name(), key), md5.New())
		if err != nil {
			b.Errorf(" %s", err)
		}
		if checksums[key]["md5"] != fileChk {
			b.Error("Expected", checksums[key]["md5"], "but returned", fileChk)
		}
	}

}

// Results running all on a single thread with Payload.Add happening inside
// the Walkfunc.
// go test -bench . -benchmem -benchtime 10m
// BEFORE refactor, running as a single function.
// BenchmarkPayload	5000000000	         7.13 ns/op	       0 B/op	       0 allocs/op
// BenchmarkPayload	10000000000	         3.43 ns/op	       0 B/op	       0 allocs/op

// AFTER refactor to go routines
// BenchmarkPayload	2000000000	         0.01 ns/op	       0 B/op	       0 allocs/op
