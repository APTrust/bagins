// Main package to running BagIns from the commandline.
package main

import (
	"bagins"
)

func main() {
	data := map[string]string{
		"BagIt-Version":                `A metadata element MUST consist of a label, a colon, and a value, each separated by optional whitespace.  It is RECOMMENDED that lines not exceed 79 characters in length.  Long values may be continued onto the next line by inserting a newline (LF), a carriage return (CR), or carriage return plus newline (CRLF) and indenting the next line with linear white space (spaces or tabs).`,
		"Tag-File-Character-Encodeing": "UTF-8",
	}
	tagFile := bagins.TagFile{Filepath: "tagfiles/bagit.txt", Data: data}
	tagFile.Create()
}
