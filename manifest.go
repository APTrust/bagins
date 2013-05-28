/* Manages manifest format files in Bags.  As of BagIt spec 0.97 this means
   only manifest-<algo>.txt and tagmanifest-<algo>.txt files.

   For more information see:
	 manifest: http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.1.3
	 tagmanifest: http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.2.1
*/
package bagins

import (
	"hash"
	"os"
)

type Manifest struct {
	entries map[string]*os.File // Maintain a map of Files managed in the manifest, referenced by checksum.
	file    *os.File            // Actual File for the manifest itself.
	hash    hash.Hash           // Hash to use for checksums
	name    string              // name of the file to use, i.e. 'manifest' or 'tagmanifest'
}
