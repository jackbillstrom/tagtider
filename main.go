package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"tagtider/pkg/trafikverket"
	"tagtider/pkg/utils"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var apiKey string
var apiURL string
var httpClient *http.Client
var locationSignature string
var signatureToNameMap map[string]string

func main() {
	var err error
	departure := flag.String("departure", "", "Station of departure which to fetch data for")
	flag.Parse()

	if *departure == "" {
		log.Fatal("No departure station specified. Use the -departure flag to specify one.")
		os.Exit(1)
	}

	httpClient = &http.Client{}
	apiKey, apiURL, err = trafikverket.FetchAPIInfo(httpClient, *departure)
	if err != nil {
		log.Fatal("Error while fetching API credentials:", err)
		os.Exit(1)
	}

	// Fetch stations and save to persistent storage for ease of use
	stationer, err := trafikverket.FetchTrainStations(httpClient, apiKey, apiURL)
	if err != nil {
		log.Fatal("Error fetching train stations:", err)
		os.Exit(1)
	}

	signatureToNameMap = make(map[string]string)
	for _, station := range stationer {
		signatureToNameMap[station.LocationSignature] = station.AdvertisedLocationName
	}

	// Match the input flag to available stations
	for _, station := range stationer {
		if strings.EqualFold(station.AdvertisedLocationName, *departure) {
			locationSignature = station.LocationSignature
			break
		}
	}

	if locationSignature == "" {
		log.Fatal("No station could be matched with the given:", *departure)
		os.Exit(1)
	}

	// Initialize the date and time for the request
	now := time.Now()
	currentDate := now.Format("2006-01-02T15:04:05-07:00")
	fromTime := currentDate
	toTime := now.Add(2 * time.Hour).Format("2006-01-02T15:04:05-07:00")
	requestBody := trafikverket.GenerateRequestPayload(apiKey, currentDate, fromTime, toTime, locationSignature)

	jsonStr := trafikverket.MakeAPICall(httpClient, apiURL, requestBody)
	if err != nil {
		log.Fatal("[error] An error occurred while searching:", err)
		return
	}

	// Unmarshal response using UnmarshalResponse function
	response, err := trafikverket.UnmarshalResponse([]byte(jsonStr))
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
		return
	}

	// Headers
	lastModifiedStr := response.Response.Result[0].Info.Lastmodified.AttrDatetime
	lastModifiedTime, err := time.Parse(time.RFC3339, lastModifiedStr)
	if err != nil {
		// Hantera fel vid parsing av tiden
		log.Fatal("Error parsing last modified time:", err)
	}

	// Calc time diff
	diff := now.Sub(lastModifiedTime)
	minutesAgo := int(diff.Minutes())

	var readableDiff string
	if minutesAgo < 1 {
		readableDiff = "just now"
	} else if minutesAgo == 1 {
		readableDiff = "1 minute ago"
	} else {
		readableDiff = fmt.Sprintf("%d minutes ago", minutesAgo)
	}

	infoTitle := fmt.Sprintf("Information \t [Last modified %s]", readableDiff)
	columns := []table.Column{
		{Title: "Departure", Width: 15},
		{Title: "From", Width: 12},
		{Title: "To", Width: 12},
		{Title: "Track", Width: 5},
		{Title: "Train", Width: 10},
		{Title: infoTitle, Width: 50},
	}

	// Create rows from API data
	rows := make([]table.Row, 0)
	for _, result := range response.Response.Result {
		for _, trainAnnouncement := range result.TrainAnnouncement {
			var departureTimeStr, oldTimeStr string
			adTime, _ := time.Parse(time.RFC3339, trainAnnouncement.AdvertisedTimeAtLocation)

			if trainAnnouncement.EstimatedTimeIsPreliminary {
				estimatedTime, _ := time.Parse(time.RFC3339, trainAnnouncement.EstimatedTimeAtLocation)
				oldTimeStr = utils.StrikeThrough(adTime.Format("15:04"))            // Add strike through to the old time
				departureTimeStr = estimatedTime.Format("15:04") + " " + oldTimeStr // Merge new and old time
			} else {
				departureTimeStr = adTime.Format("15:04")
			}

			// Match LocationSignature to AdvertisedLocationName
			fromLocationName := ""
			if len(trainAnnouncement.FromLocation) > 0 {
				fromLocationName = signatureToNameMap[trainAnnouncement.FromLocation[0].LocationName]
			}

			// Match LocationSignature to AdvertisedLocationName
			toLocationName := ""
			if len(trainAnnouncement.ToLocation) > 0 {
				toLocationName = signatureToNameMap[trainAnnouncement.ToLocation[0].LocationName]
			}

			// Tidy information strings
			information := ""
			if len(trainAnnouncement.Booking) > 1 {
				information = fmt.Sprintf("%s » %s | %s", trainAnnouncement.ProductInformation[0].Description, trainAnnouncement.Booking[0].Description, trainAnnouncement.Booking[1].Description)
			} else if len(trainAnnouncement.Booking) > 0 {
				information = fmt.Sprintf("%s » %s", trainAnnouncement.ProductInformation[0].Description, trainAnnouncement.Booking[0].Description)
			} else {
				information = fmt.Sprintf("%s", trainAnnouncement.ProductInformation[0].Description)
			}

			// Append row to table
			rows = append(rows, table.Row{
				departureTimeStr,
				fromLocationName,
				toLocationName,
				trainAnnouncement.TrackAtLocation,
				fmt.Sprintf("%v %s", trainAnnouncement.Operator, trainAnnouncement.AdvertisedTrainIdent),
				information,
			})
		}
	}

	// Create table
	t := table.New(table.WithColumns(columns), table.WithRows(rows))
	t.Focus() // Ge tabellen fokus direkt

	// Set styles
	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(true)
	s.Header.Foreground(lipgloss.Color("#ffffff"))
	s.Selected = s.Selected.Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#d70000")).Bold(false)
	t.SetStyles(s)

	// Start the bubble tea program
	m := model{table: t}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("Error starting program:", err)
		os.Exit(1)
	}
}

// Init initializes the bubble tea app
func (m model) Init() tea.Cmd {
	m.ticker = time.NewTicker(time.Minute)
	return tea.Batch(tea.Tick(time.Minute, func(t time.Time) tea.Msg {
		return tickMsg(t)
	}))
}

// Update handles UI
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.ticker != nil {
				m.ticker.Stop()
			}
			return m, tea.Quit
		}
	case tickMsg:
		// Update the table every minute
		return m.updateTable()
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// updateTable is responsible for fetching the same call again
func (m *model) updateTable() (tea.Model, tea.Cmd) {
	// Initialize the date and time for the request
	now := time.Now()
	currentDate := now.Format("2006-01-02T15:04:05-07:00")
	toTime := now.Add(2 * time.Hour).Format("2006-01-02T15:04:05-07:00")
	requestBody := trafikverket.GenerateRequestPayload(apiKey, currentDate, currentDate, toTime, locationSignature)

	jsonStr := trafikverket.MakeAPICall(httpClient, apiURL, requestBody)
	if jsonStr == "" {
		log.Fatal("[error] no data returned from trafikverket no update call ...")
	}

	// Unmarshal response using UnmarshalResponse function
	response, err := trafikverket.UnmarshalResponse([]byte(jsonStr))
	if err != nil {
		log.Fatal("[error] decoding JSON:", err)
	}

	// Create rows from API data
	rows := make([]table.Row, 0)
	for _, result := range response.Response.Result {
		for _, trainAnnouncement := range result.TrainAnnouncement {
			var departureTimeStr, oldTimeStr string
			adTime, _ := time.Parse(time.RFC3339, trainAnnouncement.AdvertisedTimeAtLocation)

			if trainAnnouncement.EstimatedTimeIsPreliminary {
				estimatedTime, _ := time.Parse(time.RFC3339, trainAnnouncement.EstimatedTimeAtLocation)
				oldTimeStr = utils.StrikeThrough(adTime.Format("15:04"))            // Add strike through to the old time
				departureTimeStr = estimatedTime.Format("15:04") + " " + oldTimeStr // Merge new and old time
			} else {
				departureTimeStr = adTime.Format("15:04")
			}

			fromLocationName := ""
			if len(trainAnnouncement.FromLocation) > 0 {
				fromLocationName = signatureToNameMap[trainAnnouncement.FromLocation[0].LocationName]
			}

			toLocationName := ""
			if len(trainAnnouncement.ToLocation) > 0 {
				toLocationName = signatureToNameMap[trainAnnouncement.ToLocation[0].LocationName]
			}

			// Tidy information strings
			information := ""
			if len(trainAnnouncement.Booking) > 1 {
				information = fmt.Sprintf("%s » %s | %s", trainAnnouncement.ProductInformation[0].Description, trainAnnouncement.Booking[0].Description, trainAnnouncement.Booking[1].Description)
			} else if len(trainAnnouncement.Booking) > 0 {
				information = fmt.Sprintf("%s » %s", trainAnnouncement.ProductInformation[0].Description, trainAnnouncement.Booking[0].Description)
			} else {
				information = fmt.Sprintf("%s", trainAnnouncement.ProductInformation[0].Description)
			}

			rows = append(rows, table.Row{
				departureTimeStr,
				fromLocationName,
				toLocationName,
				trainAnnouncement.TrackAtLocation,
				fmt.Sprintf("%v %s", trainAnnouncement.Operator, trainAnnouncement.AdvertisedTrainIdent),
				information,
			})
		}
	}

	// Update the table with new rows
	m.table.SetRows(rows)

	return m, nil
}

// View renders our TUI
func (m model) View() string {
	currentVersion := utils.GetCurrentVersion()

	tableStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		MarginRight(3)

	tableStr := tableStyle.Render(m.table.View())

	logoStr := utils.GenerateLogo(currentVersion)
	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#d70000")).Bold(true).Align(lipgloss.Center)
	styledLogo := logoStyle.Render(logoStr)

	return lipgloss.JoinHorizontal(lipgloss.Top, tableStr, styledLogo)
}
