package main

import (
	"github.com/charmbracelet/bubbles/table"
	"time"
)

// model is the application's data model.
type model struct {
	table  table.Model
	ticker *time.Ticker
}

// tickMsg is the message sent by the ticker.
type tickMsg time.Time
