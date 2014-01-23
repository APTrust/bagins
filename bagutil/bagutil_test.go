// checksums_test
package bagutil

import (
	"io/ioutil"
	"os"
	"testing"
)

var test_list map[string]string = map[string]string{
	"md5":    "9e107d9d372bb6826bd81d3542a419d6",
	"sha1":   "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12",
	"sha256": "d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e592",
	"sha512": "07e547d9586f6a73f73fbac0435ed76951218fb7d0c8d788a309d785436bbb642e93a252a954f23912547d1e8a3b5ed6e1bfd7097821233fa0538f3db854fee6",
	"sha224": "730e109bd7a8a32b1cb9d9a09aa2325d2430587ddbc0c38bad911525",
	"sha384": "ca737f1014a48f4c0b6dd43cb177b0afd9e5169367544c494011e3317dbf9a509cb1e5dc1e85a941bbee3d7f2afbc9b1",
}
var test_string string = "The quick brown fox jumps over the lazy dog"

func TestFileChecksum(t *testing.T) {
	testFile, _ := ioutil.TempFile("", "_GO_TESTFILECHECKSUM_")
	testFile.WriteString(test_string)
	testFile.Close()
	for key, sum := range test_list {
		hsh, err := LookupHash(key)
		if err != nil {
			t.Error(err)
		}
		actual, err := FileChecksum(testFile.Name(), hsh())
		if err != nil {
			t.Error(err)
		}
		if sum != actual {
			t.Error("Expected", sum, "but returned", actual, "when checking", key)
		}
	}
	os.Remove(testFile.Name())
}
