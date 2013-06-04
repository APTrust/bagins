package bagins

import (
	"errors"
	"fmt"
	"github.com/APTrust/bagins/bagutil"
	"os"
	"path"
	"regexp"
	"strings"
)

// Manifest represents information about a BagIt manifest file.  As of BagIt spec
// 0.97 this means only manifest-<algo>.txt and tagmanifest-<algo>.txt files.
//
// For more information see:
//   manifest: http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.1.3
//   tagmanifest: http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.2.1
type Manifest struct {
	name string            // Path to the
	Data map[string]string // Map of filepath key and checksum value for that file
	algo string            // Hash type to use for checksums and to concat to filename
}

// Returns a pointer to a new manifest or returns an error if improperly named.
func NewManifest(name string) (m *Manifest, err error) {
	m = new(Manifest)
	m.Data = make(map[string]string)
	m.name = path.Clean(name)
	m.algo, err = m.AlgoName()
	if !strings.HasSuffix(path.Base(m.name), ".txt") {
		err = fmt.Errorf("Manifest file %s does not end in .txt as required", path.Base(m.name))
	}
	return
}

// Returns the string of the algorithm name as indicated in the filename of
// of the manifest.  It is determined by parsing the filename as per the
// specification.
func (m *Manifest) AlgoName() (algo string, err error) {
	filename := path.Base(m.name)
	re, err := regexp.Compile(`(^.*\-)(.*)(\..*$)`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) < 2 {
		return "", errors.New("Unable to determine algorithm from filename!")
	}
	algo = matches[2]
	return algo, nil
}

func (m *Manifest) RunChecksums() []error {
	invalidSums := make([]error, 0)
	for key, sum := range m.Data {
		algo, _ := m.AlgoName()
		fileChecksum := bagutil.FileChecksum(key, algo)
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
	basepath := path.Dir(m.name)

	if err := os.MkdirAll(basepath, 0777); err != nil {
		return errors.New("Error creating manifest directory: " + err.Error())
	}

	// Create the tagfile.
	fileOut, err := os.Create(m.name)
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
	return path.Clean(m.name)
}
