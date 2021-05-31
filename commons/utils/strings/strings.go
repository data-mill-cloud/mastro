package strings

import (
	"regexp"
	"strings"
)

func SplitAndTrim(input string, separator string) []string {
	vals := strings.Split(input, separator)
	for i := range vals {
		vals[i] = strings.TrimSpace(vals[i])
	}
	return vals
}

func MatchPattern(input string, pattern string) bool {
	res, _ := regexp.MatchString(pattern, input)
	return res
}
