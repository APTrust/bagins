// Main package to running BagIns from the commandline.
package main

import (
	"bagins/bagutil"
	"fmt"
)

func main() {
	hashes := []string{"md5", "sha1", "sha256", "crc32"}
	for key := range hashes {
		fmt.Println(hashes[key], bagutil.FileChecksum("/Users/swt8w/Desktop/PresentationDryRun.mp4", hashes[key]))
	}
}
