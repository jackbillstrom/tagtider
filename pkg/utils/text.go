package utils

// StrikeThrough makes text ~~strikethrough~~
func StrikeThrough(text string) string {
	strikethrough := ""
	for _, character := range text {
		strikethrough += string(character) + "\u0336" // U+0336 is the unicode for strikethrough
	}
	return strikethrough
}
