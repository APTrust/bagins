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
	algo *bagutil.CheckAlgorithm
}

// Returns a pointer to a new manifest or returns an error if improperly named.
func NewManifest(pth string, chkAlgo *bagutil.CheckAlgorithm) (*Manifest, error) {
	if _, err := os.Stat(pth); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Unable to create manifest.  Path does not exist: %s", pth)
		} else {
			return nil, fmt.Errorf("Unexpected error creating manifest: %s", err)
		}
	}
	m := new(Manifest)
	m.Data = make(map[string]string)
	m.name = path.Join(pth, "manifest-"+strings.ToLower(chkAlgo.Name)+".txt")
	m.algo = chkAlgo
	return m, nil
}

func (m *Manifest) RunChecksums() []error {
	invalidSums := make([]error, 0)
	for key, sum := range m.Data {
		fileChecksum, err := bagutil.FileChecksum(key, m.algo.Hash)
		if sum == "" {
			m.Data[key] = fileChecksum
		}
		if sum != "" && sum != fileChecksum {
			invalidSums = append(invalidSums, errors.New(fmt.Sprintln("File checkum is not valid for", key, "!")))
		}
		if err != nil {
			invalidSums = append(invalidSums, err)
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

// Tries to parse the algorithm name from a manifest filename.  Returns
// an error if unable to do so.
func GetAlgoName(name string) (string, error) {
	filename := path.Base(name)
	re, err := regexp.Compile(`(^.*\-)(.*)(\..*$)`)
	if err != nil {
		return "", err
	}
	matches := re.FindStringSubmatch(filename)
	if len(matches) < 2 {
		return "", errors.New("Unable to determine algorithm from filename!")
	}
	algo := matches[2]
	return algo, nil
}
