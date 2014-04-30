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
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Represents the basic structure of a bag which is controlled by methods.
type Bag struct {
	pth      string // path to the bag
	payload  *Payload
	Manifest *Manifest
	tagfiles map[string]*TagFile // relative path in bag as key,
}

// METHODS FOR CREATING AND INITALIZING BAGS

/*
 Creates a new bag under the location directory and creates a bag root directory
 with the provided name.  Returns an error if the location does not exist or if the
 bag already exist.

 example:
		NewBag("archive/bags", "bag-34323", "sha256")
*/
func NewBag(location string, name string, hashName string) (*Bag, error) {
	// Create the bag object.
	bag := new(Bag)

	// Start with creating the directories.
	bag.pth = filepath.Join(location, name)
	err := os.Mkdir(bag.pth, 0755)
	if err != nil {
		return nil, err
	}
	defer bag.Save()

	// Init the manifests map and create the root manifest
	bag.Manifest, err = NewManifest(bag.Path(), hashName)
	if err != nil {
		return nil, err
	}

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
	bagit.Data.AddField(*NewTagField("BagIt-Version", "0.97"))
	bagit.Data.AddField(*NewTagField("Tag-File-Character-Encoding", "UTF-8"))

	return bagit, nil
}

/*
	Reads the directory provided as the root of a new bag and attemps to parse the file
	contents into payload, manifests and tagfiles.
*/
func ReadBag(pth string, tagfiles []string, manifest string) (*Bag, error) {
	// validate existance
	fi, err := os.Stat(pth)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a directory.", pth)
	}

	// Get the payload directory.
	payload, err := NewPayload(filepath.Join(pth, "data"))
	if err != nil {
		return nil, err
	}

	// Get the bag root directory.
	bag := new(Bag)
	bag.pth = pth
	bag.payload = payload
	bag.tagfiles = make(map[string]*TagFile)

	if manifest == "" {
		bag.findManifest()
		if bag.Manifest == nil {
			return nil, fmt.Errorf("Unable to parse a manifest")
		}
	} else {
		manifestPath := filepath.Join(bag.pth, manifest)
		if _, err := os.Stat(manifestPath); err != nil {
			return nil, fmt.Errorf("Manifest", manifest, "does not exist")
		}
		parsedManifest, errs := ReadManifest(manifestPath)
		if errs != nil && len(errs) > 0 {
			return nil, fmt.Errorf("Unable to parse a manifest", errs)
		} else {
			bag.Manifest = parsedManifest
		}
	}

	for _, tName := range tagfiles {
		tf, errs := ReadTagFile(filepath.Join(bag.pth, tName))
		if len(errs) != 0 {
			log.Println("While parsing tagfiles:", errs)
		}
		if tf != nil {
			bag.tagfiles[tName] = tf
		}
	}

	// TODO change this return
	return bag, nil
}

func (b *Bag) findManifest() {
	bagFiles, _ := b.ListFiles()
	for _, fName := range bagFiles {
		pth := filepath.Join(b.pth, fName)
		if b.Manifest == nil && strings.HasPrefix(fName, "manifest-") {
			b.Manifest, _ = ReadManifest(pth)
		}
	}
}

// METHODS FOR MANAGING BAG PAYLOADS

/*
  Adds a file specified by src parameter to the data directory under
  the relative path and filename provided in the dst parameter.
  example:
			err := b.AddFile("/tmp/myfile.txt", "myfile.txt")
*/
func (b *Bag) AddFile(src string, dst string) error {
	fx, err := b.payload.Add(src, dst, b.Manifest)
	if err != nil {
		return err
	}

	b.Manifest.Data[filepath.Join("data", dst)] = fx

	return err
}

// Performans a Bag.AddFile on all files found under the src location including all
// subdirectories.
// example:
//			errs := b.AddDir("/tmp/mypreservationfiles")
func (b *Bag) AddDir(src string) (errs []error) {
	data, errs := b.payload.AddAll(src, b.Manifest)

	for key := range data {
		b.Manifest.Data[filepath.Join("data", key)] = data[key]
	}
	return errs
}

// METHODS FOR MANAGING BAG TAG FILES

/*
 Adds a tagfile to the bag with the filename provided,
 creating whatever subdirectories are needed if supplied
 as part of name parameter.
 example:
			err := b.AddTagfile("baginfo.txt")
*/
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

/*
 Finds a tagfile in by its relative path to the bag root directory.
 example:
			tf, err := b.TagFile("bag-info.txt")
*/
func (b *Bag) TagFile(name string) (*TagFile, error) {
	if tf, ok := b.tagfiles[name]; ok {
		return tf, nil
	}
	return nil, fmt.Errorf("Unable to find tagfile %s", name)
}

/*
  Lists all the current tag files the bag is tracking.
*/
func (b *Bag) ListTagFiles() []string {
	names := make([]string, len(b.tagfiles))
	i := 0
	for k, _ := range b.tagfiles {
		names[i] = k
		i++
	}

	return names
}

/*
 Convienence method to return the bag-info.txt tag file if it exists.  Since
 this is optional it will not be created by default and will return an error
 if you have not defined or added it yourself via Bag.AddTagfile
*/
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

/*
 This method writes all the relevant tag and manifest files to finish off the
 bag.
*/
func (b *Bag) Save() (errs []error) {
	// Write all the manifest files.
	if err := b.Manifest.Create(); err != nil {
		errs = append(errs, err)
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

/*
 Walks the bag directory and subdirectories and returns the
 filepaths found inside and any errors skipping files in the
 payload directory.
*/
func (b *Bag) ListFiles() ([]string, error) {

	var files []string

	// WalkDir function to collect files in the bag..
	visit := func(pth string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		isPayload := strings.HasPrefix(pth, b.payload.Name())
		if !info.IsDir() || !isPayload {
			fp, err := filepath.Rel(b.Path(), pth)
			if err != nil {
				return err
			}
			if fp != "." {
				files = append(files, fp)
			}
		}
		return err
	}

	if err := filepath.Walk(b.Path(), visit); err != nil {
		return nil, err
	}

	return files, nil
}
