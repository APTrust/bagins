// tagfile_test
package bagins_test

import (
	"fmt"
	"github.com/APTrust/bagins"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Field TESTS

func TestNewField(t *testing.T) {
	exp_label := "test-field"
	exp_value := "this is my test"
	f := bagins.NewTagField(exp_label, exp_value)

	if f == nil { // This should test type but for the life of me I can't figure out how.
		t.Error("Tag Field object not returned")
	}
	if f.Label() != exp_label {
		t.Error("Tag Field label not created properly!")
	}
	if f.Value() != exp_value {
		t.Error("Tag Field value not created properly!")
	}
}

// Using a single test for set and get labels since they rely on one another.
func TestLabel(t *testing.T) {
	f := bagins.NewTagField("label", "value")
	if f.Label() != "label" {
		t.Error("Tag Field label not created properly!")
	}

	f.SetLabel("new-label")
	if f.Label() != "new-label" {
		t.Error("Tag Field label not reset properly!")
	}

	if f.Value() != "value" {
		t.Error("Tag Field value not set or retained properly when label changed!")
	}
}

// Using single test for set and get values since they rely on one another.
func TestValue(t *testing.T) {
	f := bagins.NewTagField("label", "value")
	if f.Value() != "value" {
		t.Error("Tag Field value not created properly!")
	}

	f.SetValue("new value")
	if f.Value() != "new value" {
		t.Error("Tag Field value not set or read properly!")
	}

	if f.Label() != "label" {
		t.Error("Tag Field label value not retained when value set!")
	}
}

// FieldList TESTS

func TestNewTagFieldList(t *testing.T) {
	var tfl interface{} = bagins.NewTagFieldList()
	if _, ok := tfl.(*bagins.TagFieldList); !ok {
		t.Error("TagFieldList not returned!")
	}
}

// Doing a unified test for Fields and SetFields
func TestFields(t *testing.T) {

	fl := bagins.NewTagFieldList()
	test_len := func(l int) { // DRY!
		if len(fl.Fields()) != l {
			t.Error("Expected TagField length of", l, "but", len(fl.Fields()), "was returned!")
		}
	}

	test_len(0)

	newFields := []bagins.TagField{
		*bagins.NewTagField("label1", "value1"),
		*bagins.NewTagField("label2", "value2"),
		*bagins.NewTagField("label3", "value3"),
	}
	fl.SetFields(newFields)
	test_len(3)

	for i := 0; i < 3; i++ {
		exp := fmt.Sprintf("label%d", i+1)
		act := fl.Fields()[i].Label()
		if exp != act {
			t.Error("Expected", exp, "but returned", act)
		}
	}
}

func TestAddField(t *testing.T) {
	fl := bagins.NewTagFieldList()
	exp_len := 100
	for i := 0; i < exp_len; i++ {
		tmp := strconv.Itoa(i)
		fl.AddField(*bagins.NewTagField(tmp, tmp))
	}

	if len(fl.Fields()) != exp_len {
		t.Error("Expected", exp_len, "fields but returned", len(fl.Fields()), "!")
	}

	for i, f := range fl.Fields() {
		if f.Value() != strconv.Itoa(i) {
			t.Error("Expected field value of", strconv.Itoa(i), "but returned", f.Value(), "!")
		}
	}
}

func TestRemoveField(t *testing.T) {
	fl := bagins.NewTagFieldList()
	test_len := func(l int) { // DRY again!
		if len(fl.Fields()) != l {
			t.Error("Expected TagField length of", l, "but", len(fl.Fields()), "was returned!")
		}
	}

	for i := 0; i < 100; i++ {
		tmp := strconv.Itoa(i)
		fl.AddField(*bagins.NewTagField(tmp, tmp))
	}
	test_len(100)

	// Should error if removing out of range.
	if err := fl.RemoveField(-6); err == nil {
		t.Error("Trying to remove negative index does not produce expected error!")
	}
	if err := fl.RemoveField(100); err == nil {
		t.Error("Trying to remove out of bound index does not produce expected error!")
	}
	test_len(100)

	// Remove every other one of the first 25 and test
	for i := 0; i < 50; i++ {
		if i%2 == 0 {
			fl.RemoveField(i)
		}
	}
	test_len(75)

}

// TagFile TESTS

func TestNewTagFile(t *testing.T) {
	_, err := bagins.NewTagFile("tagfile.txt")
	if err != nil {
		t.Error("Tagfile raised an error incorrectly!")
	}
	_, err = bagins.NewTagFile(".tagfile")
	if err == nil {
		t.Error("Bag tagfile name did not raise error as expected.")
	}
}

func TestReadTagFile(t *testing.T) {
	// Expected Data
	exp_list := [][]string{
		[]string{"description", strings.Repeat("test ", 40)},
		[]string{"title", "This is my title"},
		[]string{"description", strings.Repeat("more ", 80)},
	}

	// Prep the test file
	testPath := filepath.Join(os.TempDir(), "_GOTEST_READTAGFILE_bagit.txt")
	tagFile, _ := bagins.NewTagFile(testPath)
	for _, exp := range exp_list {
		tagFile.Data.AddField(*bagins.NewTagField(exp[0], exp[1]))
	}
	tagFile.Create()
	defer os.Remove(testPath)

	// Open and test parsing the file.
	tf, errs := bagins.ReadTagFile(testPath)
	for _, err := range errs {
		t.Error(err)
	}
	if len(tf.Data.Fields()) != 3 {
		t.Error("Expected 3 but returned", len(tf.Data.Fields()), "fields!")
	}

	fields := tagFile.Data.Fields()
	for idx, exp := range exp_list {
		if fields[idx].Label() != exp[0] {
			t.Error("Tag field", idx, "label", fields[idx].Label(), "is not expected value of", exp[0])
		}
		if fields[idx].Value() != exp[1] {
			t.Error("Tag field", idx, "value", fields[idx].Value(), "is not expected value of", exp[1])
		}
	}
}

func TestTagFileCreate(t *testing.T) {
	testPath := filepath.Join(os.TempDir(), "golang_test_tagfiles/_GOTEST_bagit.txt")
	tagFile, _ := bagins.NewTagFile(testPath)
	tagFile.Data.AddField(*bagins.NewTagField("BagIt-Version", "A metadata element MUST consist of a label, a colon, and a value, each separated by optional whitespace.  It is RECOMMENDED that lines not exceed 79 characters in length.  Long values may be continued onto the next line by inserting a newline (LF), a carriage return (CR), or carriage return plus newline (CRLF) and indenting the next line with linear white space (spaces or tabs)."))
	tagFile.Data.AddField(*bagins.NewTagField("Tag-File-Character-Encodeing", "UTF-8"))

	err := tagFile.Create()
	if err != nil {
		t.Error(err)
	}
	if _, err := os.Stat(testPath); err != nil {
		t.Error("File and path", testPath, "not created!")
	}
	os.RemoveAll(filepath.Dir(testPath))
}
