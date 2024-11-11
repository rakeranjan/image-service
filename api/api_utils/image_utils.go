package apiutils

import (
	"fmt"
	"regexp"
)

func ValidImageFile(fileName string) bool {
	regex := `(?i)^.+\.(jpg|jpeg)$`
	matched, err := regexp.MatchString(regex, fileName)
	if err != nil {
		fmt.Printf("Error validating file name: %v\n", err)
		return false
	}
	return matched
}
