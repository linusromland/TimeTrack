package utils

import (
	"fmt"
	"slices"
	"strings"
)

func Confirm(message string) bool {
	fmt.Println(message + " (Y/N)")
	var confirmation string
	_, err := fmt.Scanln(&confirmation)
	if err != nil {
		return false
	}

	fmt.Println()
	var allowedConfirmations = []string{"yes", "y"}

	return slices.Contains(allowedConfirmations, strings.ToLower(confirmation))
}
