/* Manages manifest format files in Bags.  As of BagIt spec 0.97 this means
   only manifest-<algo>.txt and tagmanifest-<algo>.txt files.

   For more information see:
	 manifest: http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.1.3
	 tagmanifest: http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.2.1
*/
package bagins

import (
	"bagins/bagutil"
	"errors"
	"fmt"
)

type Manifest struct {
	Data     map[string]string // Map of filepath key and checksum value for that file
	Filepath string            // Actual File for the manifest itself.
	Algo     string            // Hash to use for checksums
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
