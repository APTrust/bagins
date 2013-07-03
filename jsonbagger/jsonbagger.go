package jsonbagger

import (
	"github.com/APTrust/bagins"
	"github.com/APTrust/bagins/bagutil"
	"path/filepath"
	"fmt"
	"strings"
)

type TargetInfo struct {
	Dirs  []string
	Files []string
}

type BagArgs struct {
	Name     string
	Algo     string
	Targets  *TargetInfo
	TagFiles map[string]map[string] // filename:  field: value, field: value
}

// Object to register with JSON RPC server for bag creation.
type JSONBagger struct {
	dir string // directory to create bags in.
}

func (jb *JSONBagger) Create(ba *BagArgs, result *string) error {

	bag_errs := []error

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
	for tagpath, fields := range ba.TagFiles {
		tf, err := bag.AddTagfile(tagpath)
		if err != nil {
			return err
		}
		for key, value := range fields {
			tf.Data[key] = value
		}
	}

	// Process the target files.
	for _, dir := range ba.Targets.Dirs {
		errs := bag.AddDir(dir)
		if len(errs) > 0 {
			bag_errs = append(bag_errs, errs)
		}
	}
	for _, fPath := range ba.Targets.Files {
		err := bag.AddFile(fPath, filepath.Base(fPath))
		if err != nil {
			bag_errs = append(bag_errs, err)
		}
	}
	feedback := []string{fmt.Sprint("Bag Created:", bag.Path())}
	if len(bag_errs) > 0 {
		feedback = append(feedback, "Errors Reported:")
		feedback = append(feedback, bag_errs)
	}
	result = strings.Join(feedback, "\n")
	return nil
}
