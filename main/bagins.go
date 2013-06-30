// Main package to running BagIns from the commandline.
package main

import (
	"flag"
	"fmt"
	// "github.com/APTrust/bagins"
	// "github.com/APTrust/bagins/bagutil"
)

var algo string
var baseDir string
var bagName string
var srcDir string
var help bool
var opts map[string]*option

type option struct {
	full string
	df   string
	help string
}

func newOption(full string, df string, help string) *option {
	opt := new(option)
	opt.full = full
	opt.df = df
	opt.help = help
	return opt
}

func init() {
	opts["algo"] = newOption("algo", "md5", "Hash type to use for checksums.")
	opts["dir"] = newOption("dir", "", "Destination directory for the bag.")
	opts["name"] = newOption("name", "", "Name of the bag root dirctory.")
	opts["src"] = newOption("src", "", "Directory containing the files to bag.")

	flag.StringVar(&algo, "algo", opt["algo"]["default"], opt["algo"]["help"])
	flag.StringVar(&baseDir, "dir", opt["dir"]["default"], opt["dir"]["help"])
	flag.StringVar(&bagName, "name", "", "Name of the bag root directory.")
	flag.BoolVar(&help, "h", false, "Show help for this tool.")
}

func usage() {
	usg := `Usage: 
				go run bagins.go -dir dirctory -name name -filedir directory [-algo algo]

			Flags:
				`
	fmt.Println(usg)
	for _, opt := range opts {
		fmt.Println("	", "-"+opt["name"], opt["help"])
	}
}

func main() {
	flag.Parse()
	usage()

	// BASIC CODE BELOW.
	// cs, err := bagutil.NewCheckByName(algo)
	// if err != nil {
	// 	fmt.Println("Unable to find checksum", algo)
	// 	return
	// }
	// bag, err := bagins.NewBag(`C:\tmp`, "bag-all-pictures", cs)
	// if err != nil {
	// 	fmt.Println("Bag Error:", err)
	// 	return
	// }
	// errs := bag.AddDir(`E:\Pictures`)
	// for err := range errs {
	// 	fmt.Println("AddDir Error:", err)
	// 	return
	// }
	// bag.Close()
	// fmt.Println("Done!")
	// return
}
