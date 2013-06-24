// Package for dealing with bag structures.
package bagins

import (
	"errors"
	"fmt"
	"github.com/APTrust/bagins/bagutil"
	"os"
	"path"
)

// Basic type referencing main elements of a bag.
type Bag struct {
	pth       string // the bag is under.
	payload   *Payload
	manifests map[string]*Manifest
	tagfiles  map[string]*TagFile
	cs        *bagutil.ChecksumAlgorithm
}

// Creates a new bag in the provided location and name.  Returns an error
// if the location does not exist or if the bag does already exist.
func NewBag(location string, name string, cs *bagutil.ChecksumAlgorithm) (*Bag, error) {
	// Start with creating the directories.
	bagPath := path.Join(location, name)
	err := os.Mkdir(bagPath, 0755)
	if err != nil {
		return nil, err
	}
	err = os.Mkdir(path.Join(bagPath, "/data/"), 0755)
	if err != nil {
		return nil, err
	}

	// Create the bag object.
	bag := new(Bag)
	bag.cs = cs
	bag.manifests = make(map[string]*Manifest)
	bag.tagfiles = make(map[string]*TagFile)
	bag.pth = bagPath
	bag.payload, err = NewPayload(location)
	if err != nil {
		return nil, err
	}
	tf, err := bag.createBagIt()
	if err != nil {
		return nil, err
	}
	bag.tagfiles["bagit"] = tf

	return bag, nil
}

// Creates the required bagit.txt file as per the specification
// http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.1.1
func (b *Bag) createBagIt() (*TagFile, error) {
	bagit, err := NewTagFile(path.Join(b.Path(), "bagit.txt"))
	if err != nil {
		return nil, err
	}
	bagit.Data["BagIt-Version"] = "0.97"
	bagit.Data["Tag-File-Character-Encoding"] = "UTF-8"
	bagit.Create()
	return bagit, nil
}

// Adds a file to the bag payload and adds the generated checksum to the
// manifest.
func (b *Bag) PackFile(src string, dst string) error {
	return errors.New("Not implemented")
}

// Performans a Bag.Add on all files found under the src location including all
// subdirectories.
func (b *Bag) PackDir(src string) error {
	return errors.New("Not implemented")
}

func (b *Bag) AddManifest(algo string) error {
	return errors.New("Not implemented")
}

func (b *Bag) AddTagfile(name string) error {
	return errors.New("Not implemented")
}

// Returns the data fields for the baginfo.txt tag file in key, value pairs.
func (b *Bag) BagInfo() (map[string]string, error) {
	tf, err := b.TagData("bag-info")
	if err != nil {
		return nil, err
	}
	return tf, nil
}

func (b *Bag) TagData(name string) (map[string]string, error) {
	if tf, ok := b.tagfiles[name]; ok {
		return tf.Data, nil
	}
	return nil, fmt.Errorf("Unable to find tagfile %s", name)
}

func (b *Bag) ManifestData(algo string) (map[string]string, error) {
	if mf, ok := b.manifests[algo]; ok {
		return mf.Data, nil
	}
	return nil, fmt.Errorf("Unable to find manifest-%s.txt", algo)
}

func (b *Bag) Path() string {
	return b.pth
}
