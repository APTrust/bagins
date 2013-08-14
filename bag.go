/*
Package for working with files stored using the BagIt specification (see below).

It facilitates the creation of bags, adding files to the bag payload and managing
checksums for the file manifest as well as data stored in tag files.

For more information on Bag tagfiles see
http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.3
*/
package bagins

/*

“He that breaks a thing to find out what it is has left the path of wisdom.”

- Gandalf the Grey

*/

import (
	"fmt"
	"github.com/APTrust/bagins/bagutil"
	"os"
	"path/filepath"
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
//		cs := bagutil.NewChecksumAlgorithm("sha256", hsh)
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

// Adds a file specified by src parameter to the data directory under
// the relative path and filename provided in the dst parameter.
// example:
//			err := b.AddFile("/tmp/myfile.txt", "myfile.txt")
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

// Performans a Bag.AddFile on all files found under the src location including all
// subdirectories.
// example:
//			errs := b.AddDir("/tmp/mypreservationfiles")
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
// example:
// 			mf, err := b.Manifest()
func (b *Bag) Manifest() (*Manifest, error) {
	if mf, ok := b.manifests[b.cs.Name()]; ok {
		return mf, nil
	}
	return nil, fmt.Errorf("Unable to find manifest-%s.txt", b.cs.Name())
}

// METHODS FOR MANAGING BAG TAG FILES

// Adds a tagfile to the bag with the filename provided,
// creating whatever subdirectories are needed if supplied
// as part of name parameter.
// example:
// 			err := b.AddTagfile("baginfo.txt")
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

// Finds a tagfile in by its relative path to the bag root directory.
// example:
//			tf, err := b.TagFile("bag-info.txt")
func (b *Bag) TagFile(name string) (*TagFile, error) {
	if tf, ok := b.tagfiles[name]; ok {
		return tf, nil
	}
	return nil, fmt.Errorf("Unable to find tagfile %s", name)
}

// Convienence method to return the bag-info.txt tag file if it exists.  Since
// this is optional it will not be created by default and will return an error
// if you have not defined or added it yourself via Bag.AddTagfile
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

// Walks the bag directory and subdirectories and returns the
// filepaths found inside and any errors.
func (b *Bag) Contents() ([]string, []error) {

	fList := []string{}
	eList := []error{}

	// WalkDir function to collect files in the bag..
	visit := func(pth string, info os.FileInfo, err error) error {
		if err != nil {
			eList = append(eList, err)
		}
		if !info.IsDir() {
			fp, err := filepath.Rel(b.Path(), pth)
			if err != nil {
				return err
			}
			fList = append(fList, fp)
		}
		return err
	}

	if err := filepath.Walk(b.Path(), visit); err != nil {
		eList = append(eList, err)
	}

	return fList, eList
}

// Returns all the filepaths for all files being tracked by the bag.
// This includes the list of manifests, tags and files in the data directory.
// TODO:  Remove the error slice
func (b *Bag) FileManifest() ([]string, []error) {

	fList := []string{}
	eList := []error{}

	for fPath, _ := range b.tagfiles {
		fList = append(fList, fPath)
	}

	mf, _ := b.Manifest()
	for fPath, _ := range mf.Data {
		fList = append(fList, fPath)
	}
	// Confirm Manifest files are there.
	for _, mf := range b.manifests {
		fPath := filepath.Base(mf.Name())
		fList = append(fList, fPath)
	}

	return fList, eList
}

// Checks that the bag actually contains all the files it expect to and returns
// slice of errors indicating the ones that don't.
func (b *Bag) Inventory() []error {
	// Confirm Tagfiles are there.

	fls, errs := b.FileManifest()
	if len(errs) > 0 {
		return errs
	}

	for _, fl := range fls {
		if _, err := os.Stat(filepath.Join(b.Path(), fl)); os.IsNotExist(err) {
			errs = append(errs, fmt.Errorf("Unable to find: %v", fl))
		}
	}

	return errs
}

// Method returns the filepath of any files appearing in Bag.Contents that are
// not found in Bag.FileManifest
func (b *Bag) Orphans() []string {

	oList := []string{}
	// Make map to compare contents to.
	mf, _ := b.FileManifest()
	fMap := make(map[string]bool)
	for _, fn := range mf {
		fMap[fn] = true
	}

	cn, _ := b.Contents()
	for _, fn := range cn {
		if _, ok := fMap[fn]; !ok {
			oList = append(oList, fn)
		}
	}

	return oList
}
