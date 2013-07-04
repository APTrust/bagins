/*


For more information on Bag tagfiles see
http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.3
*/
package bagins

import (
	"fmt"
	"github.com/APTrust/bagins/bagutil"
	"os"
	"path/filepath"
	"strings"
)

// Represents the basic structure of a bag which is controlled by methods.
type Bag struct {
	pth       string // the bag is under.
	payload   *Payload
	manifests map[string]*Manifest // Algo string as key.
	tagfiles  map[string]*TagFile  // relative path to bag as key,
	cs        *bagutil.ChecksumAlgorithm
}

// METHODS FOR CREATING AND INITALIZING BAGS

// Creates a new bag under the location directory and creates a bag root directory
// with the provided name.  Returns an error if the location does not exist or if the
// bag already exist.
//
// example:
//		hsh, _ := bagutil.LookupHashFunc("sha256")
//		cs := bagutil.NewChecksumAlgorithm(algo, hsh)
// 		NewBag("archive/bags", "bag-34323", cs)
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

	// Init the manifests map and create the root manifest
	mf, err := NewManifest(bag.Path(), cs)
	if err != nil {
		return nil, err
	}
	bag.manifests[cs.Name()] = mf

	// Init the payload directory and such.
	plPath := filepath.Join(bag.Path(), "data")
	err = os.Mkdir(plPath, 0755)
	if err != nil {
		return nil, err
	}
	bag.payload, err = NewPayload(plPath)
	if err != nil {
		return nil, err
	}

	// Init tagfiles map and create the BagIt.txt Tagfile
	bag.tagfiles = make(map[string]*TagFile)
	tf, err := bag.createBagItFile()
	if err != nil {
		return nil, err
	}
	bag.tagfiles["bagit.txt"] = tf

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

// METHODS FOR MANAGING BAG PAYLOADS

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

// Returns the default manifest of the bag as determined by its
// checksum algorithm.
func (b *Bag) Manifest() (*Manifest, error) {
	if mf, ok := b.manifests[b.cs.Name()]; ok {
		return mf, nil
	}
	return nil, fmt.Errorf("Unable to find manifest-%s.txt", b.cs.Name())
}

// METHODS FOR MANAGING BAG TAG FILES

// Adds a tagfile to the bag, creating whatever subdirectories are needed
// as indicated by the filename.
func (b *Bag) AddTagfile(name string) error {
	tagPath := filepath.Join(b.Path(), name)
	if err := os.MkdirAll(filepath.Dir(tagPath), 0766); err != nil {
		return err
	}
	tf, err := NewTagFile(tagPath)
	if err != nil {
		return err
	}
	b.tagfiles[name] = tf
	if err := tf.Create(); err != nil {
		return err
	}
	return nil
}

// Finds a tagfile in by it's relative path to the bag root directory.
func (b *Bag) TagFile(name string) (*TagFile, error) {
	if tf, ok := b.tagfiles[name]; ok {
		return tf, nil
	}
	return nil, fmt.Errorf("Unable to find tagfile %s", name)
}

// Returns the data fields for the baginfo.txt tag file in key, value pairs.
func (b *Bag) BagInfo() (*TagFile, error) {
	tf, err := b.TagFile("bag-info.txt")
	if err != nil {
		return nil, err
	}
	return tf, nil
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
	for _, mf := range b.manifests {
		if err := mf.Create(); err != nil {
			errs = append(errs, err)
		}
	}

	// TODO Write all the tag files.
	for _, tf := range b.tagfiles {
		if err := os.MkdirAll(filepath.Dir(tf.Name()), 0766); err != nil {
			errs = append(errs, err)
		}
		if err := tf.Create(); err != nil {
			errs = append(errs, err)
		}
	}
	return
}

// Method looks to confirm that all the expected files are present in the bag.  Note
// this DOES NOT validate the bag or check the data values are accurate.  It
// only confirms the expected file structure.
func (b *Bag) Inventory() error {
	// Confirm Tagfiles are there.

	for fPath, _ := range b.tagfiles {
		if _, err := os.Stat(filepath.Join(b.Path(), fPath)); os.IsNotExist(err) {
			return fmt.Errorf("Tagfile not found: %s", fPath)
		}
	}
	// Confirm Payload files are there.
	mf, _ := b.Manifest()
	for fPath, _ := range mf.Data {
		if _, err := os.Stat(filepath.Join(b.Path(), fPath)); os.IsNotExist(err) {
			return fmt.Errorf("Payload file not found: %s", fPath)
		}
	}
	// Confirm Manifest files are there.
	for _, mf := range b.manifests {
		fPath := filepath.Base(mf.Name())
		if _, err := os.Stat(filepath.Join(b.Path(), fPath)); os.IsNotExist(err) {
			return fmt.Errorf("Manifest file not found: %s", fPath)
		}
	}

	return nil
}

// Method looks to find any items inside the bag that are not accounted for in one of
// as part of the tags, manifests or in the payload list.
func (b *Bag) Orphans() []string {
	fList := make(map[string]bool)
	oList := []string{}

	// Tag files
	for fPath, _ := range b.tagfiles {
		fList[fPath] = true
	}

	// Manifest Payload
	mf, _ := b.Manifest()
	for fPath, _ := range mf.Data {
		fList[fPath] = true
	}

	// Manifest Files
	for _, mf := range b.manifests {
		fPath := filepath.Base(mf.Name())
		fList[fPath] = true
	}

	// WalkDir function to collect files in the bag..
	visit := func(pth string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			relPath := strings.TrimPrefix(pth, b.Path()+bagutil.PathSeparator())
			if _, ok := fList[relPath]; !ok {
				oList = append(oList, relPath)
			}
		}
		return err
	}
	filepath.Walk(b.Path(), visit)

	return oList
}
