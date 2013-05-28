// Main package to running BagIns from the commandline.
package main

import (
"fmt"
"bagins"
)

func main() {
	bag := bagins.NewBag(nil, nil)
	fmt.Println(bag)
}