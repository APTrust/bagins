package jsonbagger

import (
	"errors"
	"fmt"
	"github.com/APTrust/bagins"
	"github.com/APTrust/bagins/bagutil"
	"log"
	"path/filepath"
)

type TargetInfo struct {
	Dirs  []string
	Files []string
}

type BagArgs struct {
	Name     string
	Algo     string
	Targets  *TargetInfo
	TagFiles map[string]map[string]string
}

// Object to register with JSON RPC server for bag creation.
type JSONBagger struct {
	dir string // directory to create bags in.
} 

func NewJSONBagger(dir string) *JSONBagger {
	jb := new(JSONBagger)
	jb.dir = dir
	return jb
}

func (jb *JSONBagger) Create(ba *BagArgs, result *string) error {

	bag_errs := []error{}

	// Setup the checksum algorithm.
	cs, err := bagutil.NewCheckByName(ba.Algo)
	if err != nil {
		return err
	}

	// Setup the bag itself.
	bag, err := bagins.NewBag(jb.dir, ba.Name, cs)
	if err != nil {
		return err
	}

	// Write the tag files and field values.
	for tagpath, fieldMap := range ba.TagFiles {
		err := bag.AddTagfile(tagpath)
		if err != nil {
			return err
		}
		tf, _ := bag.TagFile(tagpath)
		for key, value := range fieldMap {
			tf.Data[key] = value
		}
	}

	// Process the target files.
	for _, dir := range ba.Targets.Dirs {
		copy(bag_errs, bag.AddDir(dir))
	}
	for _, fPath := range ba.Targets.Files {
		err := bag.AddFile(fPath, filepath.Base(fPath))
		if err != nil {
			bag_errs = append(bag_errs, err)
		}
	}
	if err := bag.Close(); err != nil {
		log.Println(err)
	}
	*result = fmt.Sprint("Bag Created:", bag.Path())

	if len(bag_errs) > 0 {
		for _, err := range bag_errs {
			log.Println(err)
		}
		return errors.New("Errors encountered when creating bag, see the log for information.")
	}
	return nil
}
