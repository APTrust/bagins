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

type TagFile struct {
	name string            // Filepath for tag file.
	Data map[string]string // key value pairs of data for the tagfile.
}

func NewTagFile(name string) (tf *TagFile, err error) {
	_, err = os.Stat(filepath.Dir(name))
	re, _ := regexp.Compile(`.*\.txt`)
	if !re.MatchString(filepath.Base(name)) {
		err = errors.New(fmt.Sprint("Tagfiles must end in .txt and contain at least 1 letter.  Provided: ", filepath.Base(name)))
	}
	tf = new(TagFile)
	tf.name = filepath.Clean(name)
	tf.Data = make(map[string]string)
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
