package utils

const Version = "1.0.0"

// GetCurrentVersion returns the current version of the program. Ugly and should be refactored.
func GetCurrentVersion() string {
	return Version + "\n" + "Author: jackbillstrom @ github.com \n"
}
