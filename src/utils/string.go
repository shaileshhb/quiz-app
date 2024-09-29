package utils

import "regexp"

// ValidateString validates if the input matches the given pattern
func ValidateString(input, pattern string) (bool, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}
	return re.MatchString(input), nil
}
