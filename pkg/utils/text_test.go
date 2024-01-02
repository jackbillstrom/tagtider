package utils_test

import (
	"strings"
	"tagtider/pkg/utils"
	"testing"
)

// TestStrikeThrough tests the StrikeThrough function
func TestStrikeThrough(t *testing.T) {
	text := "test"
	strikethrough := utils.StrikeThrough(text)

	for _, char := range text {
		if !strings.Contains(strikethrough, string(char)+"\u0336") {
			t.Errorf("StrikeThrough() = %v, want it to contain %v", strikethrough, string(char)+"\u0336")
		}
	}
}
