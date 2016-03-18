package main

//
//
// This application can be compiled and deployed as a stand alone executable to
// create very basic bags from the commandline.
//
//

import (
	"flag"
	"fmt"
	"github.com/APTrust/bagins"
	"strings"
	"time"
)

var (
	dir          string
	name         string
	payload      string
	algo         string
	tagmanifests string
)

func init() {
	flag.StringVar(&dir, "dir", "", "Directory to create the bag.")
	flag.StringVar(&name, "name", "", "Name for the bag root directory.")
	flag.StringVar(&payload, "payload", "", "Directory of files to parse into the bag")
	flag.StringVar(&algo, "algo", "md5", "Checksum algorithm to use.  md5, sha1, sha224, sha256, sha512, sha384")
	flag.StringVar(&tagmanifests, "tagmanifests", "", "Set to true to create tag manifests. Default is false.")

	flag.Parse()
}

func usage() {

	usage := `
Usage: ./bagmaker -dir <value> -name <value> -payload <value> [-algo <value>]

Flags:

    -algo <value>
     Checksum algorithm to use.  md5, sha1, sha224, sha256,
     sha512, or sha384. Defaults to md5. Use commas (without
	 spaces) to specify multiple algorithms. E.g. md5,sha256

    -dir <value>
     Directory to create the bag.

    -name <value>
     Name for the bag root directory.

    -payload <value>
     Directory of files to copy into the bag.

    -tagmanifests <value>
     Set to true to create tag manifests. Default is false.

     Example:

     Put all of /home/joe into a bag called joes_bag in the current working
     directory. Create payload manifests with md5 and sha256 checksums, and
     create tagmanifests as well (using the same checksum algorithms).

     bagmaker -payload /home/joe -name joes_bag -dir . -algo md5,sha256 -tagmanifests true

`
	fmt.Println(usage)
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

	algoList := parseAlgorithms(algo)

	createTagManifests := false
	if tagmanifests == "true" {
		createTagManifests = true
	}

	begin := time.Now()

	bag, err := bagins.NewBag(dir, name, algoList, createTagManifests)
	if err != nil {
		fmt.Println("Bag Error:", err)
		return
	}

	errs := bag.AddDir(payload)
	for idx := range errs {
		fmt.Println("AddDir Error:", errs[idx])
		return
	}

	bag.Save()

	elapsed := time.Since(begin)
	fmt.Println("END: elapsed in", elapsed.Seconds(), "seconds.")
	return
}

// Parses command line arguments to go into the

// func parse_info(args []string) map[string]string {
// 	info := make(map[string]string)
// 	return info
// }

func parseAlgorithms(algo string) (algorithms []string) {
	if algo == "" {
		algorithms = []string { "md5" }
	} else {
		algorithms = strings.Split(algo, ",")
	}
	return algorithms
}
