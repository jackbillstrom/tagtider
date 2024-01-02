package utils_test

import (
	"strings"
	"tagtider/pkg/utils"
	"testing"
)

// TestGenerateLogo tests the GenerateLogo function
func TestGenerateLogo(t *testing.T) {
	version := "1.0.0"
	logo := utils.GenerateLogo(version)

	if !strings.Contains(logo, version) {
		t.Errorf("GenerateLogo() = %v, want it to contain %v", logo, version)
	}
}
