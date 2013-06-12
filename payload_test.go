package bagins_test

import (
	"github.com/APTrust/bagins"
	//"io/ioutil"
	//"crypto/md5"
	"io/ioutil"
	"os"
	"path"
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

// func TestPayloadAdd(t *testing.T) {
// 	pDir, _ := ioutil.TempDir("", "_GOTEST")

// 	p, err := bagins.NewPayload(pDir)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	testFile, _ := ioutil.TempFile("", "_GO_TESTFILE_")
// 	testFile.WriteString("Test the checksum")
// 	testFile.Close()

// 	chkSum, err := p.Add(testFile.Name(), path.Base(testFile.Name(), md5.New()))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	exp := "92d7a9f0f4a30ca782dcae5fe83ca7eb"
// 	if exp != chkSum {
// 		t.Error("Checksum", chkSum, "did not match", exp)
// 	}
// 	os.RemoveAll(pDir)
// }
