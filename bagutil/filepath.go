package bagutil

/*

"Where will wants not, a way opens."

- Eowyn

*/

import (
	"os"
)

// Utility method to return the operation system seperator as a string.
func PathSeparator() string {
	return string(byte(os.PathSeparator))
}
