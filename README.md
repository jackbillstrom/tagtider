# Tågtider
A simple CLI tool for fetching train departures from Trafikverket's API.

![Screenshot of the application in use with Ubuntu mono as font in the Standard terminal](screenshot.png)

## Requirements
- [Go](https://golang.org/dl/)
- A terminal emulator that supports [ANSI escape codes](https://en.wikipedia.org/wiki/ANSI_escape_code#Colors). Any modern terminal emulator should work.

## Installation
The easiest way to install is to grab the latest release from [releases]() and run the following command:


**Linux**
```bash
wget https://github.com/jackbillstrom/tagtider/releases/download/v1.0.0/tagtider-linux-amd64 -O tagtider && chmod +x tagtider && sudo mv tagtider /usr/local/bin/

# Or via manual steps
wget https://github.com/jackbillstrom/tagtider/releases/download/v1.0.0/tagtider-linux-amd64 -O tagtider
chmod +x tagtider
sudo mv tagtider /usr/local/bin/
```

**MacOS**
```bash
# Download the binary
curl -Lo tagtider https://github.com/jackbillstrom/tagtider/releases/download/v1.0.0/tagtider-darwin-amd64

# Make it executable
chmod +x tagtider

# Move it to a directory in your PATH (optional)
sudo mv tagtider /usr/local/bin/
```

**Windows**
```powershell
# Download the binary
Invoke-WebRequest -Uri https://github.com/jackbillstrom/tagtider/releases/download/v1.0.0/tagtider-windows-amd64.exe -OutFile tagtider.exe

# Move it to a directory in your PATH (optional)
Move-Item tagtider.exe C:\bin\
```
## Usage
```bash
tagtider --departure <station>
```

Example of asking for departures from Stockholm Central Station:
```bash
tagtider --departure "Stockholm C"
```

## License
[Gnu GPL v3](https://www.gnu.org/licenses/gpl-3.0.en.html)

## Disclaimer
This project is not affiliated with Trafikverket in any way. It is simply a hobby project, and should be treated as such.

## TODO
- [ ] Add tabs for departure and arrival
- [ ] Highlight cancelled trains
- [ ] Enable autocomplete for stations?
- [ ] Write more tests

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## Acknowledgements
- [Trafikverket](https://www.trafikverket.se/)
- [Bubbletea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)
