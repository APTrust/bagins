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
	defer bag.Close()
	bag.pth = bagPath
	bag.cs = cs
	bag.manifests = make(map[string]*Manifest)
	bag.manifests[cs.Name()], err = NewManifest(bag.Path(), cs)
	if err != nil {
		return nil, err
	}
	bag.tagfiles = make(map[string]*TagFile)

	bag.payload, err = NewPayload(location)
	if err != nil {
		return nil, err
	}
	tf, err := bag.createBagItFile()
	if err != nil {
		return nil, err
	}
	bag.tagfiles["bagit"] = tf

	// TODO initiate a baginfo.txt file as well, even if it's blank.

	return bag, nil
}

// Creates the required bagit.txt file as per the specification
// http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.1.1
func (b *Bag) createBagItFile() (*TagFile, error) {
	bagit, err := NewTagFile(path.Join(b.Path(), "bagit.txt"))
	if err != nil {
		return nil, err
	}
	bagit.Data["BagIt-Version"] = "0.97"
	bagit.Data["Tag-File-Character-Encoding"] = "UTF-8"
	return bagit, nil
}

// Adds a file to the bag payload and adds the generated checksum to the
// manifest.
func (b *Bag) PackFile(src string, dst string) error {
	return errors.New("Not implemented")
}

// Performans a Bag.Add on all files found under the src location including all
// subdirectories.
func (b *Bag) PackDir(src string) (errs []error) {
	data, errs := b.payload.AddAll(src, b.cs.Algo())
	if mf, err := b.Manifest(); err != nil {
		errs = append(errs, err)
	} else {
		for key := range data {
			mf.Data[key] = data[key]
		}
	}

	return errs
}

func (b *Bag) AddManifest(algo string) error {
	return errors.New("Not implemented")
}

func (b *Bag) AddTagfile(name string) error {
	return errors.New("Not implemented")
}

// Returns the data fields for the baginfo.txt tag file in key, value pairs.
func (b *Bag) BagInfo() (*TagFile, error) {
	tf, err := b.TagFile(b.Path() + "bag-info")
	if err != nil {
		return nil, err
	}
	return tf, nil
}

func (b *Bag) TagFile(name string) (*TagFile, error) {
	if tf, ok := b.tagfiles[name]; ok {
		return tf, nil
	}
	return nil, fmt.Errorf("Unable to find tagfile %s", name)
}

func (b *Bag) Manifest() (*Manifest, error) {
	mf, err := b.ManifestFile(b.cs.Name())
	if err != nil {
		return nil, err
	}
	return mf, nil
}

func (b *Bag) ManifestFile(algo string) (*Manifest, error) {
	if mf, ok := b.manifests[algo]; ok {
		return mf, nil
	}
	return nil, fmt.Errorf("Unable to find manifest-%s.txt", algo)
}

func (b *Bag) Path() string {
	return b.pth
}

// This method writes all the relevant tag and manifest files to finish off the
// bag.
func (b *Bag) Close() (errs []error) {

	// Write all the manifest files.
	for key := range b.manifests {
		if mf, err := b.ManifestFile(key); err != nil {
			errs = append(errs, err)
		} else {
			if err = mf.Create(); err != nil {
				errs = append(errs, err)
			}
		}

	}

	// TODO Write all the tag files.
	for key := range b.tagfiles {
		if tf, err := b.TagFile(key); err != nil {
			errs = append(errs, err)
		} else {
			if err = tf.Create(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}
