package utils

// TrafikverketLogo is the Trafikverket logo in Braille
var TrafikverketLogo = `
⠀⠀⠀⠀⢐⡀⠀⠀⠀⠀
⠀⡠⢀⠄⢀⡁⠠⡐⢄⠀
⢫⠀⡮⠀⢸⡅⠀⡵⠀⡕
⠀⢤⢤⣢⢤⡤⣤⢬⠤⠀
⠀⠀⠉⠈⣗⣇⠉⠁⠁⠀
⢫⠆⠀⠀⣗⡵⠀⠀⠰⡽
⠨⡣⠀⠀⠑⠋⠀⠀⢜⠅
⠀⠘⠪⠖⡦⣔⠦⠗⠁⠀
`

// GenerateLogo generates the Trafikverket logo with the current version
func GenerateLogo(version string) string {
	return TrafikverketLogo + "\nVersion: " + version
}
