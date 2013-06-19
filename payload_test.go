package bagins_test

import (
	"crypto/md5"
	"github.com/APTrust/bagins"
	"github.com/APTrust/bagins/bagutil"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestNewPayload(t *testing.T) {

	tmpPyld := path.Join(os.TempDir(), "__GOTEST_Payload/")

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
	pDir, _ := ioutil.TempDir("", "_GOTEST_Payload")
	defer os.RemoveAll(pDir)

	p, err := bagins.NewPayload(pDir)
	if err != nil {
		t.Error(err)
	}

	testFile, _ := ioutil.TempFile("", "_GO_TESTFILE_")
	testFile.WriteString("Test the checksum")
	testFile.Close()
	defer os.Remove(testFile.Name())

	chkSum, err := p.Add(testFile.Name(), path.Base(testFile.Name()), md5.New())
	if err != nil {
		t.Error(err)
	}
	exp := "92d7a9f0f4a30ca782dcae5fe83ca7eb"
	if exp != chkSum {
		t.Error("Checksum", chkSum, "did not match", exp)
	}
}

func TestPayloadAddAll(t *testing.T) {
	// Make src temp dir
	srcDir, _ := ioutil.TempDir("", "_GOTEST_SRCDIR_")
	defer os.RemoveAll(srcDir)
	pDir, _ := ioutil.TempDir("", "_GOTEST_Payload_")
	defer os.RemoveAll(pDir)

	// Make src temp test files
	for i := 0; i < 100; i++ {
		tstFile, _ := ioutil.TempFile(srcDir, "_GOTEST_FILE_")
		tstFile.WriteString("Test the checksum")
		tstFile.Close()
	}

	p, _ := bagins.NewPayload(pDir)
	fxs, errs := p.AddAll(srcDir, md5.New())
	if errs != nil {
		t.Errorf("Add all returned %d errors", len(errs))
	}
	for key := range fxs {
		fileChk, err := bagutil.FileChecksum(path.Join(p.Name(), key), md5.New())
		if err != nil {
			t.Errorf(" %s", err)
		}
		if fxs[key] != fileChk {
			t.Error("Expected", fxs[key], "but returned", fileChk)
		}
	}

}

func BenchmarkPayload(b *testing.B) {
	srcDir, _ := ioutil.TempDir("", "_GOTEST_SRCDIR_")
	defer os.RemoveAll(srcDir)
	pDir, _ := ioutil.TempDir("", "_GOTEST_Payload_")
	defer os.RemoveAll(pDir)

	// Make src temp test files
	for i := 0; i < 5000; i++ {
		tstFile, _ := ioutil.TempFile(srcDir, "_GOTEST_FILE_")
		tstFile.WriteString(strings.Repeat("Test the checksum. ", 100000))
		tstFile.Close()
	}

	b.ResetTimer()

	p, _ := bagins.NewPayload(pDir)

	fxs, err := p.AddAll(srcDir, md5.New())
	if err != nil {
		b.Error(err)
	}

	b.StopTimer()

	// Make sure the actual values check out.
	for key := range fxs {
		fileChk, err := bagutil.FileChecksum(path.Join(p.Name(), key), md5.New())
		if err != nil {
			b.Errorf(" %s", err)
		}
		if fxs[key] != fileChk {
			b.Error("Expected", fxs[key], "but returned", fileChk)
		}
	}

}
