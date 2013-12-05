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
	"github.com/APTrust/bagins/bagutil"
	"os"
	"time"
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

	usage := `Usage:
	./bagmaker -dir <value> -name <value> -payload <value> [-algo <value>]

Flags:

	-algo <value> Checksum algorithm to use.  md5, sha1, sha224, sha256, 
	              sha512, or sha384. Defaults to md5.

	-dir <value> Directory to create the bag.

	-name <value> Name for the bag root directory.

	-payload <value> Directory of files to parse into the bag.

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

	begin := time.Now()

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
	for idx := range errs {
		fmt.Println("AddDir Error:", errs[idx])
		return
	}

	// if info := parse_info(os.Args); len(info) > 0 {
	// 	bag.AddTagfile("bag-info.txt")
	// 	if tf, err := bag.BagInfo(); err != nil {
	// 		tf.Data = info
	// 	}
	// }

	bag.Close()

	elapsed := time.Since(begin)
	fmt.Println("END: elapsed in", elapsed.Seconds(), "seconds.")
	return
}

// Parses command line arguments to go into the

func parse_info(args []string) map[string]string {
	info := make(map[string]string)
	return info
}
