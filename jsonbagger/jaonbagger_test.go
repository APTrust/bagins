package jsonbagger_test

import (
	"github.com/APTrust/bagins/jsonbagger"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func setUpBagAttributes(srcDir string) *jsonbagger.BagArgs {
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

	return ba
}

func setUpTargetFiles(num int) string {
	// Setup the Test files
	srcDir, _ := ioutil.TempDir("", "_GOTEST_PAYLOAD_SRC_")
	for i := 0; i < num; i++ {
		fi, _ := ioutil.TempFile(srcDir, "TEST_GO_ADDFILE_")
		fi.WriteString("Test the checksum")
		fi.Close()
	}
	return srcDir
}

func TestCreate(t *testing.T) {

	srcDir := setUpTargetFiles(50)
	defer os.RemoveAll(srcDir)

	ba := setUpBagAttributes(srcDir)

	// Setup the bagger itself.
	jb := jsonbagger.NewJSONBagger(os.TempDir())
	var result string
	err := jb.Create(ba, &result)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(filepath.Join(os.TempDir(), ba.Name))
}
