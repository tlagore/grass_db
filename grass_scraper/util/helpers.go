package util

import (
	"os"
	"strings"
)

func GetCharAt(s string, ch byte) int {
	i := -1

	for index, _ := range s {
		if s[index] == ch {
			return index
		}
	}

	return i
}

func DirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return false, err
}

func Contains(lst []string, s string, ignoreCase bool) bool {
	if ignoreCase {
		s = strings.ToLower(s)
	}

	for _, val := range lst {
		var copyVal string

		if ignoreCase {
			copyVal = strings.ToLower(val)
		}else{
			copyVal = val
		}

		if copyVal == s {
			return true
		}
	}

	return false
}