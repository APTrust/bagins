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
	"os"
	"path/filepath"
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
	algo *bagutil.ChecksumAlgorithm
}

// Returns a pointer to a new manifest or returns an error if improperly named.
func NewManifest(pth string, chkAlgo *bagutil.ChecksumAlgorithm) (*Manifest, error) {
	if _, err := os.Stat(pth); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Unable to create manifest.  Path does not exist: %s", pth)
		} else {
			return nil, fmt.Errorf("Unexpected error creating manifest: %s", err)
		}
	}
	m := new(Manifest)
	m.Data = make(map[string]string)
	m.name = filepath.Join(pth, "manifest-"+strings.ToLower(chkAlgo.Name())+".txt")
	m.algo = chkAlgo
	return m, nil
}

func ReadManifest(name string) (*Manifest, []error) {
	var errs []error

	algo, err := GetAlgoName(name)
	if err != nil {
		return nil, append(errs, err)
	}

	hsh, err := bagutil.NewCheckByName(algo)
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

	m, err := NewManifest(name, hsh)
	if err != nil {
		return nil, append(errs, err)
	}
	m.Data = data

	return m, errs

}

func (m *Manifest) RunChecksums() []error {
	invalidSums := make([]error, 0)
	for key, sum := range m.Data {
		fileChecksum, err := bagutil.FileChecksum(key, m.algo.New())
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

// Returns a sting of the filename for this manifest file based on Path, BaseName and Algo
func (m *Manifest) Name() string {
	return filepath.Clean(m.name)
}

// Tries to parse the algorithm name from a manifest filename.  Returns
// an error if unable to do so.
func GetAlgoName(name string) (string, error) {
	filename := filepath.Base(name)
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

// Reads the contents of file and parses checksum and file information in manifest format as
// per the bagit specification.
func parseManifestData(file *os.File) (map[string]string, []error) {
	var errs []error
	// See regexp examples at http://play.golang.org/p/_msLJ-lBEu
	re := regexp.MustCompile(`^(\S*) *(.\S*)`)

	scanner := bufio.NewScanner(file)
	values := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			data := re.FindStringSubmatch(line)
			values[data[1]] = data[2]
		} else {
			errs = append(errs, fmt.Errorf("Unable to parse data from line: %s", line))
		}
	}

	return values, errs
}
