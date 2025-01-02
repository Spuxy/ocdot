// Package foundation f
package foundation

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ClearErrorMsg type is
type ClearErrorMsg struct{}

// ClearErrorAfter func
func ClearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return ClearErrorMsg{}
	})
}
