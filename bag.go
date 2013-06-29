// Package for dealing with bag structures.
package bagins

import (
	"fmt"
	"github.com/APTrust/bagins/bagutil"
	"os"
	"path/filepath"
)

// Basic type referencing main elements of a bag.
type Bag struct {
	pth       string // the bag is under.
	payload   *Payload
	manifests map[string]*Manifest
	tagfiles  map[string]*TagFile
	cs        *bagutil.ChecksumAlgorithm
}

// METHODS FOR CREATING AND INITALIZING BAGS

// Creates a new bag in the provided location and name.  Returns an error
// if the location does not exist or if the bag does already exist.
func NewBag(location string, name string, cs *bagutil.ChecksumAlgorithm) (*Bag, error) {
	// Start with creating the directories.
	bagPath := filepath.Join(location, name)
	err := os.Mkdir(bagPath, 0755)
	if err != nil {
		return nil, err
	}

	// Create the bag object.
	bag := new(Bag)
	defer bag.Close()
	bag.pth = bagPath
	bag.cs = cs
	bag.manifests = make(map[string]*Manifest)

	if err = bag.AddManifest(cs.Name()); err != nil {
		return nil, err
	}
	bag.tagfiles = make(map[string]*TagFile)

	// Make the payload directory and such.
	plPath := filepath.Join(bag.Path(), "data")
	err = os.Mkdir(plPath, 0755)
	if err != nil {
		return nil, err
	}
	bag.payload, err = NewPayload(plPath)
	if err != nil {
		return nil, err
	}

	// Create the BagIt.txt Tagfile
	tf, err := bag.createBagItFile()
	if err != nil {
		return nil, err
	}
	bag.tagfiles["bagit"] = tf

	return bag, nil
}

// Creates the required bagit.txt file as per the specification
// http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.1.1
func (b *Bag) createBagItFile() (*TagFile, error) {
	if err := b.AddTagfile("bagit.txt"); err != nil {
		return nil, err
	}
	bagit, err := b.TagFile("bagit.txt")
	if err != nil {
		return nil, err
	}
	bagit.Data["BagIt-Version"] = "0.97"
	bagit.Data["Tag-File-Character-Encoding"] = "UTF-8"
	return bagit, nil
}

// METHODS FOR MANAGING BAG PAYLOAD

// Adds a file to the bag payload and adds the generated checksum to the
// manifest.
func (b *Bag) AddFile(src string, dst string) error {
	fx, err := b.payload.Add(src, dst, b.cs.New())
	if err != nil {
		return err
	}
	if mf, err := b.Manifest(); err == nil {
		mf.Data[filepath.Join("data", dst)] = fx
	}
	return err
}

// Performans a Bag.Add on all files found under the src location including all
// subdirectories.
func (b *Bag) AddDir(src string) (errs []error) {
	data, errs := b.payload.AddAll(src, b.cs.Algo())
	mf, err := b.Manifest()
	if err != nil {
		errs = append(errs, err)
	}
	for key := range data {
		mf.Data[filepath.Join("data", key)] = data[key]
	}
	return errs
}

// METHODS FOR MANAGING BAG MANIFESTS

func (b *Bag) AddManifest(algo string) error {
	hsh, err := bagutil.LookupHashFunc(algo)
	if err != nil {
		return err
	}
	cs := bagutil.NewChecksumAlgorithm(algo, hsh)
	if err != nil {
		return err
	}
	mf, err := NewManifest(b.Path(), cs)
	if err != nil {
		return err
	}
	b.manifests[algo] = mf
	return nil
}

func (b *Bag) ManifestFile(algo string) (*Manifest, error) {
	if mf, ok := b.manifests[algo]; ok {
		return mf, nil
	}
	return nil, fmt.Errorf("Unable to find manifest-%s.txt", algo)
}

func (b *Bag) Manifest() (*Manifest, error) {
	mf, err := b.ManifestFile(b.cs.Name())
	if err != nil {
		return nil, err
	}
	return mf, nil
}

// METHODS FOR MANAGING BAG TAG FILES

func (b *Bag) AddTagfile(name string) error {
	tf, err := NewTagFile(filepath.Join(b.Path(), name))
	if tf != nil {
		b.tagfiles[name] = tf
	}
	return err
}

// Returns the data fields for the baginfo.txt tag file in key, value pairs.
func (b *Bag) BagInfo() (*TagFile, error) {
	tf, err := b.TagFile("bag-info.txt")
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

// TODO create methods for managing fetch file.

// TODO create methods to manage tagmanifest files.

// METHODS FOR MANAGING OR RETURNING INFORMATION ABOUT THE BAG ITSELF

// Returns the full path of the bag including it's own directory.
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

// TODO create a method to return the name of the bag root folder alone as the
// bag name

// TODO create method to return a list of tag files.

// TODO create a method to return a list of manifest files.
