// Main package to running BagIns from the commandline.
package main

import (
	//"bagins"
	"bagins/bagutil"
	"fmt"
	"os"
)

func main() {
	//bag := bagins.NewBag(nil, nil)
	//fmt.Println(bag)

	file, err := os.Open("/Users/swt8w/Desktop/PresentationDryRun.mp4")
	if err != nil {
		panic(err)
	}
	fmt.Println(bagutil.Sha1Checksum(file))

}
