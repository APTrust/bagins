// Main package to running BagIns from the commandline.
package main

import (
	"flag"
	"fmt"
	"github.com/APTrust/bagins"
	"github.com/APTrust/bagins/bagutil"
)

var (
	dir     string
	name    string
	payload string
	algo    string
)

func init() {
	flag.StringVar(&dir, "dir", "", "Directory to create the bag.")
	flag.StringVar(&name, "name", "", "Name for the bag root directory.")
	flag.StringVar(&payload, "payload", "", "Directory of files to parse into the bag")
	flag.StringVar(&algo, "algo", "md5", "Checksum algorithm to use.  md5, sha1, sha224, sha256, sha512, sha384")

	flag.Parse()
}

func usage() {
	fmt.Println(`
Usage:
	go run baging.go -dir <value> -name <value> -payload <value> [-algo <value>]

Flags:
		`)
	printFlags := func(a *flag.Flag) {
		fmt.Println("	-"+a.Name+" <value>", a.Usage)
	}

	flag.VisitAll(printFlags)
}

func main() {

	if dir == "" {
		usage()
		return
	}
	if name == "" {
		usage()
		return
	}
	if payload == "" {
		usage()
		return
	}

	cs, err := bagutil.NewCheckByName(algo)
	if err != nil {
		fmt.Println("Unable to find checksum", algo)
		return
	}
	bag, err := bagins.NewBag(dir, name, cs)
	if err != nil {
		fmt.Println("Bag Error:", err)
		return
	}
	errs := bag.AddDir(payload)
	for err := range errs {
		fmt.Println("AddDir Error:", err)
		return
	}
	bag.Close()
	fmt.Println("Done!")
	return
}
