// Main package to running BagIns from the commandline.
package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/APTrust/bagins"
)

func main() {
	srcDir := "/Users/swt8w/Documents/"

	dstDir := "/Users/swt8w/payload_test"

	p, _ := bagins.NewPayload(dstDir)
	fxs, errs := p.AddAll(srcDir, sha256.New())
	if errs != nil {
		fmt.Println(errs)
	}

	for key := range fxs {
		fmt.Println(key, fxs[key])
	}

	fmt.Println("Done procesing", len(fxs), "files.")

}
