package bagins

/*

"Oft in lies truth is hidden."

- Glorfindel

*/

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/APTrust/bagins/bagutil"
	"hash"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

/*
 Manifest represents information about a BagIt manifest file.  As of BagIt spec
 0.97 this means only manifest-<algo>.txt and tagmanifest-<algo>.txt files.

 For more information see:
   manifest: http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.1.3
   tagmanifest: http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.2.1
*/
type Manifest struct {
	name          string            // Path to the manifest file
	manifestType  string            // payload manifest or tag manifest?
	Data          map[string]string // Key is file path, value is checksum
	hashName      string
	hashFunc      func() hash.Hash
}

const (
	PayloadManifest = "payload_manifest"
	TagManifest     = "tag_manifest"
)

// Returns a pointer to a new manifest or returns an error if improperly named.
func NewManifest(pathToFile string, hashName string, manifestType string) (*Manifest, error) {
	if manifestType != PayloadManifest && manifestType != TagManifest {
		return nil, fmt.Errorf("Param manifestType must be either bagins.PayloadManifest " +
			"or bagins.TagManifest")
	}
	if _, err := os.Stat(filepath.Dir(pathToFile)); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Unable to create manifest. Path does not exist: %s", pathToFile)
		} else {
			return nil, fmt.Errorf("Unexpected error creating manifest: %s", err)
		}
	}
	m := new(Manifest)
	m.hashName = strings.ToLower(hashName)
	hashFunc, err := bagutil.LookupHash(hashName)
	if err != nil {
		return nil, err
	}
	m.hashFunc = hashFunc
	m.Data = make(map[string]string)

	// Older versions allow pathToFile to be empty...
	if !strings.HasSuffix(pathToFile, "manifest-" + hashName + ".txt") {
		if manifestType == PayloadManifest {
			pathToFile = filepath.Join(pathToFile, "manifest-" + m.hashName +".txt")
		} else {
			pathToFile = filepath.Join(pathToFile, "tagmanifest-" + m.hashName + ".txt")
		}
	}

	m.name = pathToFile
	m.manifestType = PayloadManifest
	if manifestType == TagManifest {
		m.manifestType = TagManifest
	}

	return m, nil
}

/*
  Opens a manifest file, parses attemps to parse the hashtype from the filename, parses
  the file contents and returns a pointer to a Manifest.  Error slice may comprise multiple
  parsing errors when attempting to read data for fault tolerance.
*/
func ReadManifest(name string) (*Manifest, []error) {
	var errs []error

	hashName, err := parseAlgoName(name)
	if err != nil {
		return nil, append(errs, err)
	}

	file, err := os.Open(name)
	if err != nil {
		return nil, append(errs, err)
	}

	data, e := parseManifestData(file)
	if e != nil {
		errs = append(errs, e...)
	}

	manifestType := PayloadManifest
	if strings.HasPrefix(path.Base(name), "tagmanifest-") {
		manifestType = TagManifest
	}
	m, err := NewManifest(name, hashName, manifestType)
	if err != nil {
		return nil, append(errs, err)
	}
	m.Data = data

	return m, errs

}

/*
  Calculates a checksum for files listed in the manifest and compares it to the value
  stored in manifest file.  Returns an error for each file that fails the fixity check.
*/
func (m *Manifest) RunChecksums() []error {
	var invalidSums []error

	for key, sum := range m.Data {
		pathToFile := filepath.Join(filepath.Dir(m.name), key)
		fileChecksum, err := bagutil.FileChecksum(pathToFile, m.hashFunc())
		if sum != fileChecksum {
			invalidSums = append(invalidSums, fmt.Errorf("File checksum %s is not valid for %s:%s", sum, key, fileChecksum))
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
	basepath := filepath.Dir(m.name)

	if err := os.MkdirAll(basepath, 0777); err != nil {
		return err
	}

	// Create the tagfile.
	fileOut, err := os.Create(m.name)
	if err != nil {
		return err
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

// Returns the contents of the manifest in the form of a string.
// Useful if you don't want to write directly to disk.
func (m *Manifest) ToString() string {
	str := ""
	for fName, ckSum := range m.Data {
		str += fmt.Sprintf("%s %s\n", ckSum, fName)
	}
	return str
}


// Returns a sting of the filename for this manifest file based on Path, BaseName and Algo
func (m *Manifest) Name() string {
	return filepath.Clean(m.name)
}

// Returns the name of the manifest's hashing algorithm.
// "sha256", "md5", etc.
func (m *Manifest) Algorithm() string {
	return m.hashName
}

// Returns the type of manifest. Either 'payload' or 'tag'.
func (m *Manifest) Type() string {
	return m.manifestType
}

// Tries to parse the algorithm name from a manifest filename.  Returns
// an error if unable to do so.
func parseAlgoName(name string) (string, error) {
	filename := filepath.Base(name)
	re, err := regexp.Compile(`(^.*\-)(.*)(\.txt$)`)
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

// Reads the contents of file and parses checksum and file information in manifest format as
// per the bagit specification.
func parseManifestData(file *os.File) (map[string]string, []error) {
	var errs []error
	// See regexp examples at http://play.golang.org/p/_msLJ-lBEu
	// Regex matches these reqs from the bagit spec: "One or
	// more linear whitespace characters (spaces or tabs) MUST separate
	// CHECKSUM from FILENAME." as specified here:
	// http://tools.ietf.org/html/draft-kunze-bagit-10#section-2.1.3
	re := regexp.MustCompile(`^(\S*)\s*(.*)`)

	scanner := bufio.NewScanner(file)
	values := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			data := re.FindStringSubmatch(line)
			values[data[2]] = data[1]
		} else {
			errs = append(errs, fmt.Errorf("Unable to parse data from line: %s", line))
		}
	}

	return values, errs
}
