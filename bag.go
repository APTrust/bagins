// Package for dealing with bag structures.
package bagins

import (
	"fmt"
	"os"
	"path"
)

// Basic type referencing main elements of a bag.
type Bag struct {
	name     string   // Name of the bag, will also be the top level directory name.
	path     string   // the bag is under.
	manifest Manifest // Required manifest file
	data     *os.File // Data Directory
	bagit    TagFile  // bagit.txt Tag file.
}

// Creates a new bag in the provided location and name.  Returns an error
// if the location does not exist or if the bag does already exist.
func NewBag(location string, name string) (*Bag, error) {
	baseDir := path.Clean(location)
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("Bag destination location does not exist! Returned: %v", err)
	}
	if _, err := os.Stat(path.Join(baseDir, name)); os.IsExist(err) {
		return nil, fmt.Errorf("Bag %s already exists!", path.Join(baseDir, name))
	}
	return new(Bag), nil
}
