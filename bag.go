// Package for dealing with bag structures.
package bagins

import (
	"errors"
	"fmt"
	"os"
	"path"
)

// Basic type referencing main elements of a bag.
type Bag struct {
	pth      string    // the bag is under.
	payload *Payload
	manifests map[string]Manifest
	tagfiles map[string]Tagfile
}

// Creates a new bag in the provided location and name.  Returns an error
// if the location does not exist or if the bag does already exist.
func NewBag(location string, name string) (*Bag, error) {
	baseDir := path.Clean(location)
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("Destination path does not exist! Returned: %v", err)
	}
	if _, err := os.Stat(path.Join(baseDir, name)); os.IsExist(err) {
		return nil, fmt.Errorf("Bag %s already exists! Returned: %s", path.Join(baseDir, name), err)
	}
	bagPath := os.Join(locaiton, name)
	bag := new(Bag)
	// make the bag dir
	// make the baginfo.txt file
	// make the payload directory and initialize the data dir.
	// make manifest.
	return bag, nil
}

// Adds a file to the bag payload and adds the generated checksum to the
// manifest.
func (b *Bag) AddFile(src string, dst string) error {
	return errors.New("Not implemented")
}

// Performans a Bag.Add on all files found under the src location including all
// subdirectories.
func (b *Bag) AddDir(src string) error {
	return errors.New("Not implemented")
}

func (b *Bag) AddManifest(algo string) error {
	return errors.New("Not implemented")
}

func (b *Bag) AddTagfile(name string) error {
	return errros.New("Not implemented")
}

// Returns the data fields for the baginfo.txt tag file in key, value pairs.
func (b *Bag) BagInfo() map[string]string, error {
	return b.TagFile("bag-info")
}

func (b *Bag) TagData(name string) map[string]string, error {
	data, ok := b.tagfiles[name]; ok {
		return data, nil
	}
	return nil, fmt.Errorf("Unable to find tagfile %s", name)
}

func (b *Bag) ManifestData(algo string) map[string]string, error {
	data, ok := b.manifests[algo]; ok {
		return data, nil
	}
}