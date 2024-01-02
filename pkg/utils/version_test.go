package utils_test

import (
	"strings"
	"tagtider/pkg/utils"
	"testing"
)

// TestGenerateLogo tests the GenerateLogo function
func TestGetCurrentVersion(t *testing.T) {
	version := utils.Version
	author := "Author: jackbillstrom @ github.com"
	currentVersion := utils.GetCurrentVersion()

	if !strings.Contains(currentVersion, version) || !strings.Contains(currentVersion, author) {
		t.Errorf("GetCurrentVersion() = %v, want it to contain %v and %v", currentVersion, version, author)
	}
}
