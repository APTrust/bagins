package bagins

import (
	"bagins/bagutil"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// Manifest represents information about a BagIt manifest file.  As of BagIt spec
// 0.97 this means only manifest-<algo>.txt and tagmanifest-<algo>.txt files.
//
// For more information see:
//   manifest: http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.1.3
//   tagmanifest: http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.2.1
type Manifest struct {
	Data     map[string]string // Map of filepath key and checksum value for that file
	Path     string            // Path to the Manifest File
	BaseName string            // Base filename for manifest to combine with Algo for filename.
	Algo     string            // Hash type to use for checksums and to concat to filename
}

func (m *Manifest) RunChecksums() []error {
	invalidSums := make([]error, 0)
	for key, sum := range m.Data {
		fileChecksum := bagutil.FileChecksum(key, m.Algo)
		if sum == "" {
			m.Data[key] = fileChecksum
		}
		if sum != "" && sum != fileChecksum {
			invalidSums = append(invalidSums, errors.New(fmt.Sprintln("File checkum is not valid for", key, "!")))
		}
	}
	return invalidSums
}

// Writes key value pairs to a manifest file.
func (m *Manifest) Create() error {
	if m.Name() == "" {
		return errors.New("Manifest must have values for basename and algo set to create a file.")
	}
	// Create directory if needed.
	basepath := path.Dir(m.Path)

	// kind of weird but Join is more performant.
	filename := strings.Join([]string{m.BaseName, "-", m.Algo, ".txt"}, "")

	if err := os.MkdirAll(basepath, 0777); err != nil {
		return errors.New("Error creating manifest directory: " + err.Error())
	}

	// Create the tagfile.
	fileOut, err := os.Create(path.Join(basepath, filename))
	if err != nil {
		return errors.New("Error creating manifest file: " + err.Error())
	}
	defer fileOut.Close()

	// Write fields and data to the file.
	for fName, ckSum := range m.Data {
		_, err := fmt.Fprintln(fileOut, ckSum, fName)
		if err != nil {
			return errors.New("Error writing line to manifest: " + err.Error())
		}
	}
	return nil
}

// Returns a sting of the filename for this manifest file based on Path, BaseName and Algo
func (m *Manifest) Name() string {
	if m.BaseName == "" || m.Algo == "" {
		return ""
	}
	filename := strings.Join([]string{m.BaseName, "-", m.Algo, ".txt"}, "")
	return path.Join(path.Dir(m.Path), filename)
}
