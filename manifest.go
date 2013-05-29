/* Manages manifest format files in Bags.  As of BagIt spec 0.97 this means
   only manifest-<algo>.txt and tagmanifest-<algo>.txt files.

   For more information see:
	 manifest: http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.1.3
	 tagmanifest: http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.2.1
*/
package bagins

type Manifest struct {
	Data     map[string]string // Map of file checksum and filepath
	Filepath string            // Actual File for the manifest itself.
	Algo     string            // Hash to use for checksums
}
