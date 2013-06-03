package bagins_test

import (
	"bagins"
	"os"
	"path"
	"testing"
)

func TestNewBag(t *testing.T) {
	badLocation := path.Join(os.TempDir(), "/GOTESTNOT_EXISTs/")
	_, err := bagins.NewBag(badLocation, "_GOTESTBAG_")
	if err == nil {
		t.Error("NewBag function does not recognize when a directory does not exist!")
	}

}
