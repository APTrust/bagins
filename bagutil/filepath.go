package bagutil

import (
	"os"
)

func PathSeparator() string {
	return string(byte(os.PathSeparator))
}
