package utils

import (
	"fmt"
	"slices"
	"strings"
)

func Confirm(message string) bool {
	fmt.Println(message + " (Y/N)")
	var confirmation string
	fmt.Scanln(&confirmation)
	fmt.Println()
	var allowedConfirmations = []string{"yes", "y"}

	return slices.Contains(allowedConfirmations, strings.ToLower(confirmation))
}