// Package for dealing with bag structures.
package bagins

// Basic type referencing main elements of a bag.
type Bag struct {
	Name string // Name of the bag
}

// Creates a new bag in the provided location and name.
func NewBag(location string, name string) *Bag {
	// TODO Check that location exists.
	// TODO Make a directory with name in location.
	// TODO Make a data diretory in the named bag.
	// TODO Initialize a bag-info.txt file.
	return new(Bag)
}
