package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type model struct {
	filepicker   filepicker.Model
	selectedFile string
	quitting     bool
	err          error
}
type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selectedFile = path
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}
func (m model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		s.WriteString("Pick a file:")
	} else {
		s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}

var rootCmd = &cobra.Command{
	Use:   "ocdot",
	Short: "ocdot is dotfile managment",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()

		if err != nil {
			log.Printf("Error occured: %s", err.Error())
			os.Exit(1)
		}
		var dotfilesPath = fmt.Sprintf("%s/%s", homeDir, pathToDotfiles)

		if err := checkIfStowExists(); err != nil {
			log.Printf("Error occured: %s", err.Error())
			os.Exit(1)
		}

		fp := filepicker.New()
		fp.AllowedTypes = []string{".mod", ".sum", ".go", ".txt", ".md"}
		fp.CurrentDirectory = dotfilesPath

		m := model{
			filepicker: fp,
		}
		tm, _ := tea.NewProgram(&m).Run()
		mm := tm.(model)
		fmt.Println("\n  You selected: " + m.filepicker.Styles.Selected.Render(mm.selectedFile) + "\n")
	},
}

func checkIfStowExists() error {
	var isExists bool

	if _, err := os.Stat("/usr/bin/stow"); err == nil {
		fmt.Println("lol")
		isExists = true
	}

	if _, err := os.Stat("/opt/homebrew/bin/stow"); err == nil {
		isExists = true
	}

	if !isExists {
		return errors.New("stow binary is not installed")
	}

	return nil
}

var pathToDotfiles string

func main() {
	rootCmd.PersistentFlags().StringVar(&pathToDotfiles, "path-to-file", ".dotfiles", "path to dotfile (default is $HOME/.dotfiles)")
	if err := rootCmd.Execute(); err != nil {
		log.Printf("Error occured: %s", err.Error())
		os.Exit(1)
	}
}
