package util

import (
	"bytes"
	"strings"
)

func Join(sep string, vals ...string) string {
	return strings.Join(vals, sep)
}

func RemoveDupOf(s string, target rune) string {
	var buf bytes.Buffer
	var last rune
	for i, r := range s {
		if (i == 0 || r != target) ||
			(r == target && r != last) {

			buf.WriteRune(r)
			last = r
		}
	}
	return buf.String()
}
