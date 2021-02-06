package main

import "strings"

type Case int

const (
	Lower Case = iota
	Upper
	IgnoreCase
)

func getCharAt(s string, ch byte) int {
	i := -1

	for index, _ := range s {
		if s[index] == ch {
			return index
		}
	}

	return i
}

func contains(lst []string, s string, ignoreCase bool) bool {
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