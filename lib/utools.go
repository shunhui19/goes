package lib

import "unicode"

// UcFirst ToUpper the first char of strings.
func UcFirst(str string) string {
	for _, v := range str {
		f := string(unicode.ToUpper(v))
		return f + str[len(f):]
	}
	return ""
}
