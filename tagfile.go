package bagins

/*

"Oft the unbidden guest proves the best company."

- Eomer

*/

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// TAG FIELD

/*
 Represents a tag field as referenced in the standard BagIt tag file and used
 in bag-info.txt.  It represents a standard key value pair with the label with corresponding
 value.  For more information see
 http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.2.2
*/
type TagField struct {
	label string // Name of the tag field
	value string // Value of the tag field
}

// Creates and returns a pointer to a new TagField
func NewTagField(label string, value string) *TagField {
	return &TagField{label, value}
}

// Returns the label string for the tag field.
func (f *TagField) Label() string {
	return f.label
}

// Sets the label string for the tag field.
func (f *TagField) SetLabel(l string) {
	f.label = l
}

// Returns the value string for the tag field.
func (f *TagField) Value() string {
	return f.value
}

// Sets the value string for the tag file.
func (f *TagField) SetValue(v string) {
	f.value = v
}

// TAG FIELD LIST

/*
 Represents an ordered list of tag fields as specified for use with bag-info.txt
 in the bag it standard.  It supports ordered, repeatable fields.
 http://tools.ietf.org/html/draft-kunze-bagit-09#section-2.2.2
*/
type TagFieldList struct {
	fields []TagField // Some useful manipulations in https://code.google.com/p/go-wiki/wiki/SliceTricks
}

// Returns a pointer to a new TagFieldList.
func NewTagFieldList() *TagFieldList {
	return new(TagFieldList)
}

// Returns a slice copy of the current tag fields.
func (fl *TagFieldList) Fields() []TagField {
	return fl.fields
}

// Sets the tag field slice to use for the tag field list.
func (fl *TagFieldList) SetFields(fields []TagField) {
	fl.fields = fields
}

// Adds a Field to the end of the tag field list.
func (fl *TagFieldList) AddField(field TagField) {
	fl.fields = append(fl.Fields(), field)
}

/*
 Removes a field from the tag field list at the specified index.  Returns an error if
 index out of bounds.
*/
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

// Represents a tag file object in the bag with its related fields.
type TagFile struct {
	name string        // Filepath for tag file.
	Data *TagFieldList // key value pairs of data for the tagfile.
}

/*
 Creates a new tagfile object and returns it or returns an error if improperly formatted.
 The name argument represents the filepath of the tagfile, which must end in txt
*/
func NewTagFile(name string) (tf *TagFile, err error) {
	err = validateTagFileName(name)
	tf = new(TagFile)
	tf.name = filepath.Clean(name)
	tf.Data = new(TagFieldList)
	return tf, err
}

/*
 Reads a tagfile, parsing the contents as tagfile field data and returning the TagFile object.
 name is the filepath to the tag file.  It throws an error if contents cannot be properly parsed.
*/
func ReadTagFile(name string) (*TagFile, []error) {
	var errs []error

	file, err := os.Open(name)
	if err != nil {
		return nil, append(errs, err)
	}
	defer file.Close()

	tf, err := NewTagFile(name)
	if err != nil {
		return nil, append(errs, err)
	}

	data, errs := parseTagFields(file)
	tf.Data.SetFields(data)

	return tf, errs
}

// Returns the named filepath of the tagfile.
func (tf *TagFile) Name() string {
	return tf.name
}

/*
 Creates the named tagfile and writes key value pairs to it, with indented
 formatting as indicated in the BagIt spec.
*/
func (tf *TagFile) Create() error {
	// Create directory if needed.
	if err := os.MkdirAll(filepath.Dir(tf.name), 0777); err != nil {
		return err
	}

	// Create the tagfile.
	fileOut, err := os.Create(tf.Name())
	if err != nil {
		return err
	}
	defer fileOut.Close()

	// Write fields and data to the file.
	for _, f := range tf.Data.Fields() {
		field, err := formatField(f.Label(), f.Value())
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

// Returns the contents of the tagfile in the form of a string.
// This is an alternative to Create(), which writes to disk.
func (tf *TagFile) ToString() (string, error) {
	str := ""
	for _, f := range tf.Data.Fields() {
		field, err := formatField(f.Label(), f.Value())
		if err != nil {
			return "", err
		}
		str += fmt.Sprintf("%s\n", field)
	}
	return str, nil
}


/*
Takes a tag field key and data and wraps lines at 79 with indented spaces as
per recommendation in spec.
*/
func formatField(key string, data string) (string, error) {
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

/*
 Reads the contents of file and parses tagfile fields from the contents or returns an error if
 it contains unparsable data.
*/
func parseTagFields(file *os.File) ([]TagField, []error) {
	var errors []error
	re, err := regexp.Compile(`^(\S*\:)?(\s.*)?$`)
	if err != nil {
		errors = append(errors, err)
		return nil, errors
	}

	scanner := bufio.NewScanner(file)
	var fields []TagField
	var field TagField

	// Parse the remaining lines.
	for scanner.Scan() {
		line := scanner.Text()
		// See http://play.golang.org/p/zLqvg2qo1D for some testing on the field match.
		if re.MatchString(line) {
			data := re.FindStringSubmatch(line)
			data[1] = strings.Replace(data[1], ":", "", 1)
			if data[1] != "" {
				if field.Label() != "" {
					fields = append(fields, field)
				}
				field = *NewTagField(data[1], strings.Trim(data[2], " "))
				continue
			}
			value := strings.Trim(data[2], " ")
			field.SetValue(strings.Join([]string{field.Value(), value}, " "))

		} else {
			err := fmt.Errorf("Unable to parse tag data from line: %s", line)
			errors = append(errors, err)
		}
	}
	if field.Label() != "" {
		fields = append(fields, field)
	}

	if scanner.Err() != nil {
		errors = append(errors, scanner.Err())
	}

	// See http://play.golang.org/p/nsw9zsAEPF for some testing on the field match.
	return fields, errors
}
