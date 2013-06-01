// tagfile_test
package bagins_test

import (
	"bagins"
	"os"
	"path"
	"strings"
	"testing"
)

/*
Test for Creation of the tagfile in a directory.
TODO test for formatting once those methods are written.
*/
func TestTagFileCreate(t *testing.T) {
	testPath := path.Join(os.TempDir(), "golang_test_tagfiles/_GOTEST_bagit.txt")
	tagFile, _ := bagins.NewTagFile(testPath)
	tagFile.Data = map[string]string{
		"BagIt-Version":                `A metadata element MUST consist of a label, a colon, and a value, each separated by optional whitespace.  It is RECOMMENDED that lines not exceed 79 characters in length.  Long values may be continued onto the next line by inserting a newline (LF), a carriage return (CR), or carriage return plus newline (CRLF) and indenting the next line with linear white space (spaces or tabs).`,
		"Tag-File-Character-Encodeing": "UTF-8",
	}
	tagFile.Create()
	if _, err := os.Stat(testPath); err != nil {
		t.Error("File and path", testPath, "not created!")
	}
	os.RemoveAll(path.Dir(testPath))
}

func TestFormatField(t *testing.T) {
	resultString := bagins.FormatField("tst", strings.Repeat("test ", 20))
	exp := 80
	act := strings.Index(resultString, "\n")
	if exp != act {
		t.Errorf("Found newline at %d but expected %d", act, exp)
	}
}
