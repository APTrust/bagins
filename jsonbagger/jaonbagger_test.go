package jsonbagger_test

import (
	"github.com/APTrust/bagins/jsonbagger"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCreate(t *testing.T) {

	// Setup the Test files
	srcDir, _ := ioutil.TempDir("", "_GOTEST_PAYLOAD_SRC_")
	for i := 0; i < 50; i++ {
		fi, _ := ioutil.TempFile(srcDir, "TEST_GO_ADDFILE_")
		fi.WriteString("Test the checksum")
		fi.Close()
	}
	defer os.RemoveAll(srcDir)

	// Setup the Bag Attributes
	tagfiles := make(map[string]map[string]string)
	baginfo := make(map[string]string)
	baginfo["Source-Organization"] = "APTrust"
	baginfo["Contact-Name"] = "Scott Turnbull"
	baginfo["Contact-Email"] = "scott.turnbull@aptrust.org"
	tagfiles["baginfo.txt"] = baginfo

	tgt := &jsonbagger.TargetInfo{
		Dirs:  []string{srcDir},
		Files: []string{},
	}

	ba := &jsonbagger.BagArgs{
		Name:     "_GOTEST_JSONBAGGER_CREATE_",
		Algo:     "md5",
		TagFiles: tagfiles,
		Targets:  tgt,
	}

	// Setup the bagger itself.
	jb := jsonbagger.NewJSONBagger(os.TempDir())
	var result string
	err := jb.Create(ba, &result)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(filepath.Join(os.TempDir(), ba.Name))
}
