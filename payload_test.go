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

	tmpPyld := filepath.Join(os.TempDir(), "__GOTEST_Payload/")

	// Check for failure on non-existant directory.
	_, err := bagins.NewPayload(tmpPyld)
	if err == nil {
		t.Errorf("Unexpected error return checking for non-existed directory: %s", err)
	}

	// Check for positive return when directory exists.
	pth, err := ioutil.TempDir("", "_GOTEST_Payload")
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
	os.Remove(pth)
}

func TestPayloadName(t *testing.T) {
	pDir, _ := ioutil.TempDir("", "_GOTEST_Payload_")
	defer os.Remove(pDir)

	p, _ := bagins.NewPayload(pDir)

	if pDir != p.Name() {
		t.Errorf("Payload name %s did not equal expected %s", p.Name(), pDir)
	}
}

func TestPayloadAdd(t *testing.T) {
	pDir, _ := ioutil.TempDir("", "GOTEST_Payload")
	m, _ := bagins.NewManifest(os.TempDir(), "md5")
	defer os.RemoveAll(pDir)

	p, err := bagins.NewPayload(pDir)
	if err != nil {
		t.Error(err)
	}

	testFile, _ := ioutil.TempFile("", "_GO_TESTFILE_")
	testFile.WriteString("Test the checksum")
	testFile.Close()
	defer os.Remove(testFile.Name())

	chkSum, err := p.Add(testFile.Name(), filepath.Base(testFile.Name()), m)
	if err != nil {
		t.Error(err)
	}
	exp := "92d7a9f0f4a30ca782dcae5fe83ca7eb"
	if exp != chkSum {
		t.Error("Checksum", chkSum, "did not match", exp)
	}
}

func TestPayloadAddAll(t *testing.T) {
	// Setup directories to test on
	srcDir, _ := ioutil.TempDir("", "_GOTEST_SRCDIR_")
	defer os.RemoveAll(srcDir)
	pDir, _ := ioutil.TempDir("", "_GOTEST_Payload_")
	defer os.RemoveAll(pDir)

	m, _ := bagins.NewManifest(os.TempDir(), "md5")

	// Setup test files
	for i := 0; i < 100; i++ {
		tstFile, _ := ioutil.TempFile(srcDir, "_GOTEST_FILE_")
		tstFile.WriteString("Test the checksum")
		tstFile.Close()
	}

	p, _ := bagins.NewPayload(pDir)
	fxs, errs := p.AddAll(srcDir, m)

	// It should not return an error.
	if errs != nil {
		t.Errorf("Add all returned %d errors", len(errs))
	}
	// It should have fixity values for 100 files
	if len(fxs) != 100 {
		t.Errorf("Expected 100 fixity values but returned %d", len(fxs))
	}

	for key := range fxs {
		fileChk, err := bagutil.FileChecksum(filepath.Join(p.Name(), key), md5.New())
		if err != nil {
			t.Errorf(" %s", err)
		}
		if fxs[key] != fileChk {
			t.Error("Expected", fxs[key], "but returned", fileChk)
		}
	}

}

func TestPayloadOctetStreamSum(t *testing.T) {
	// Setup Test directory
	pDir, _ := ioutil.TempDir("", "_GOTEST_Payload_")
	defer os.RemoveAll(pDir)

	// Setup test files
	for i := 0; i < 100; i++ {
		tstFile, _ := ioutil.TempFile(pDir, "_GOTEST_FILE_")
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
	srcDir, _ := ioutil.TempDir("", "_GOTEST_SRCDIR_")
	defer os.RemoveAll(srcDir)
	pDir, _ := ioutil.TempDir("", "_GOTEST_Payload_")
	defer os.RemoveAll(pDir)

	m, _ := bagins.NewManifest(os.TempDir(), "md5")

	// Make src temp test files
	for i := 0; i < 300; i++ {
		tstFile, _ := ioutil.TempFile(srcDir, "_GOTEST_FILE_")
		tstFile.WriteString(strings.Repeat("Test the checksum. ", 500000)) // produces ~9 meg text file.
		tstFile.Close()
	}

	b.ResetTimer()

	p, _ := bagins.NewPayload(pDir)

	fxs, err := p.AddAll(srcDir, m)
	if err != nil {
		b.Error(err)
	}

	b.StopTimer()

	// Make sure the actual values check out.
	for key := range fxs {
		fileChk, err := bagutil.FileChecksum(filepath.Join(p.Name(), key), md5.New())
		if err != nil {
			b.Errorf(" %s", err)
		}
		if fxs[key] != fileChk {
			b.Error("Expected", fxs[key], "but returned", fileChk)
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
