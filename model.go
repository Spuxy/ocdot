// package main is the main package.
package main

import "github.com/charmbracelet/bubbles/table"

// Model is the application's main model.
type Model struct {
	Files        []string
	SyncState    string
	SelectedFile string
	Table        table.Model
}

// ClearErrorMsg type is
type ClearErrorMsg struct{}
