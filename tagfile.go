/*
Methods for managing tag files inside of bags.  These represent the standard
tag files that are part of every Bag spec or optional tag files which can
be used liberally throughout the bag itself.

For more information on Bag tagfiles see
http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.3
*/
package bagins

import (
	"os"
)

type TagFile struct {
	file *os.File          // File this information is stored in.
	name string            // Name of the tagfile itself.
	data map[string]string // key value pairs of data for the tagfile.
}
