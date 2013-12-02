package bagins

/*

"Oft the unbidden guest proves the best company."

- Eomer

*/

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// TAG FIELD

// Represents a tag field as referenced in the standard BagIt tag file and used
// in bag-info.txt.  It represents a standard key value pair with the label with corresponding
// value.  For more information see
// http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.2.2
type TagField struct {
	label string // Name of the tag field
	value string // Value of the tag field
}

func NewTagField(label string, value string) *TagField {
	return &TagField{label, value}
}

func (f *TagField) Label() string {
	return f.label
}

func (f *TagField) SetLabel(l string) {
	f.label = l
}

func (f *TagField) Value() string {
	return f.value
}

func (f *TagField) SetValue(v string) {
	f.value = v
}

// TAG FIELD LIST

// Ordered list of fields used in a TagFile.
type TagFieldList struct {
	fields []TagField // Some useful manipulations in https://code.google.com/p/go-wiki/wiki/SliceTricks
}

func NewTagFieldList() *TagFieldList {
	return new(TagFieldList)
}

func (fl *TagFieldList) Fields() []TagField {
	return fl.fields
}

func (fl *TagFieldList) SetFields(fields []TagField) {
	fl.fields = fields
}

func (fl *TagFieldList) AddField(field TagField) {
	fl.fields = append(fl.Fields(), field)
}

func (fl *TagFieldList) RemoveField(i int) error {
	if i+1 > len(fl.Fields()) || i < 0 {
		return errors.New("Invalid index for TagField")
	}
	if len(fl.fields) == i {
		fl.fields = fl.Fields()[:i]
		return nil
	}
	fl.fields = append(fl.Fields()[:i], fl.Fields()[i+1:]...)
	return nil
}

// TAG FILES

type TagFile struct {
	name string            // Filepath for tag file.
	Data map[string]string // key value pairs of data for the tagfile.
}

// Creates a new tagfile object and returns it or returns an error if improperly formatted.
// The name argument represents the filepath of the tagfile, which must end in txt
func NewTagFile(name string) (tf *TagFile, err error) {
	err = validateTagFileName(name)
	tf = new(TagFile)
	tf.name = filepath.Clean(name)
	tf.Data = make(map[string]string)
	return tf, err
}

// Opens a tagfile, parsing the contents as tagfile field data and returning the TagFile object.
// name is the filepath to the tag file.  It throws an error if contents cannot be properly parsed.
func OpenTagFile(name string) (*TagFile, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tf, err := NewTagFile(name)
	return tf, err
}

// Returns the named filepath of the tagfile.
func (tf *TagFile) Name() string {
	return tf.name
}

// Creates the named tagfile and writes key value pairs to it, with indented
// formatting as indicated in the BagIt spec.
func (tf *TagFile) Create() error {
	// Create directory if needed.
	if err := os.MkdirAll(filepath.Dir(tf.name), 0777); err != nil {
		return err
	}

	// Create the tagfile.
	fileOut, err := os.Create(tf.name)
	if err != nil {
		return err
	}
	defer fileOut.Close()

	// Write fields and data to the file.
	for key, data := range tf.Data {
		field, err := FormatField(key, data)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(fileOut, field)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
Takes a tag field key and data and wraps lines at 79 with indented spaces as
per recommendation in spec.
*/
func FormatField(key string, data string) (string, error) {
	delimeter := "\n   "
	var buff bytes.Buffer

	// Initiate it by writing the proper key.
	writeLen, err := buff.WriteString(fmt.Sprintf("%s: ", key))
	if err != nil {
		return "", err
	}
	splitCounter := writeLen

	words := strings.Split(data, " ")

	for word := range words {
		if splitCounter+len(words[word]) > 79 {
			splitCounter, err = buff.WriteString(delimeter)
			if err != nil {
				return "", err
			}
		}
		writeLen, err = buff.WriteString(strings.Join([]string{" ", words[word]}, ""))
		if err != nil {
			return "", err
		}
		splitCounter += writeLen

	}
	return buff.String(), nil
}

// Some private convenence methods for manipulating tag files.
func validateTagFileName(name string) (err error) {
	_, err = os.Stat(filepath.Dir(name))
	re, _ := regexp.Compile(`.*\.txt`)
	if !re.MatchString(filepath.Base(name)) {
		err = errors.New(fmt.Sprint("Tagfiles must end in .txt and contain at least 1 letter.  Provided: ", filepath.Base(name)))
	}
	return err
}

// Reads the contents of file and parses tagfile fields from the contents or returns an error if
// it contains unparsable data.
// func parseFields(data map[string]string, file *os.File) (err error) {
// 	re, err := regexp.Compile(`^\S*\:`)
// 	return err
// }
