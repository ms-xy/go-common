package stringutils

import (
	"strings"
)

// Split2 is a convenience wrapper around strings.SplitN that returns 2 string
// slices no matter what the input, as one would expect SplitN to do ... rather
// silly, that something like this is needed.
func Split2(in, del string) (string, string) {
	r := strings.SplitN(in, del, 2)
	if len(r) > 1 {
		return r[0], r[1]
	} else {
		return r[0], ""
	}
}
