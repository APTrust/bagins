package bagins_test

import (
	"fmt"
	"github.com/APTrust/bagins"
	"os"
	"path"
	"testing"
)

func TestNewBag(t *testing.T) {

	// It should raise an error if the destination dir does not exist.
	badLocation := path.Join(os.TempDir(), "/GOTESTNOT_EXISTs/")
	_, err := bagins.NewBag(badLocation, "_GOFAILBAG_")
	if err == nil {
		t.Error("NewBag function does not recognize when a directory does not exist!")
	}

	// It should raise an error if the bag already exists.
	os.MkdirAll(path.Join(badLocation, "_GOFAILBAG_"), 0766)
	defer os.RemoveAll(badLocation)
	_, err = bagins.NewBag(badLocation, "_GOFAILBAG_")
	if err == nil {
		t.Error("Error not thrown when bag already exists as expected.")
	}

	// Test making an actual bag.
	bag, err := bagins.NewBag(os.TempDir(), "_GOTESTBAG_")
	fmt.Println(bag.Path())
	defer os.RemoveAll(bag.Path())
	if err != nil {
		t.Error(err)
	}
	if _, err = os.Stat(path.Join(os.TempDir(), "_GOTESTBAG_")); os.IsNotExist(err) {
		t.Errorf("Bag directory does not exist! Returned: %v", err)
	}
}
