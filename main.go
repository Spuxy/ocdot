// package ocdot
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type Model struct {
	Files        []string
	SyncState    string
	SelectedFile string
	Table        table.Model
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

// CurrentFiles func
func CurrentFiles() ([]table.Row, error) {
	var dotfilesRows []table.Row

	rootDir := "/Users/filip.boye.kofi/dot.filesbak"
	_, err := os.ReadDir(rootDir)
	if err != nil {
		fmt.Println(err)
		return []table.Row{}, err
	}
	err = filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			// fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}

		rslt := strings.Split(info.Name(), "/")
		if len(rslt) > 1 {
			return nil
		}
		// fmt.Printf("visited file or dir: %q\n", path)
		dotfilesRows = append(dotfilesRows, []string{"✅", rslt[0], "github.com/spuxy/dot.file"})
		return nil
	})

	if err != nil {
		return []table.Row{}, err
	}

	return dotfilesRows, nil
}

// Update func
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	rows, err := CurrentFiles()
	if err != nil {
		panic(err)
	}
	m.Table.SetRows(rows)
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

// Init func
func (m Model) Init() tea.Cmd {
	return nil
}

// // View Kekw
//
//	func (m Model) View() string {
//		// return fmt.Sprintf("wtf - %d", rand.Int63n(1000))
//
//		rootDir := "/Users/filip.boye.kofi/dot.filesbak"
//		_, err := os.ReadDir(rootDir)
//		if err != nil {
//			fmt.Println(err)
//			return ""
//		}
//
//		err = filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
//			if err != nil {
//				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
//				return err
//			}
//			if info.IsDir() && info.Name() == ".git" {
//				fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
//				return filepath.SkipDir
//			}
//
//			fmt.Println(strings.Split(info.Name(), "/"))
//			// fmt.Println(strings.TrimPrefix(path, rootDir))
//
//			fmt.Printf("visited file or dir: %q\n", path)
//			return nil
//		})
//		if err != nil {
//			fmt.Printf("error occurred: %v\n", err)
//			return "err"
//		}
//		return "nothing"
//	}
func (m Model) View() string {

	// m.Table.SetHeight(len(dotfilesRows) + 1)
	m.Table.SetWidth(300)

	return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).Render(m.Table.View()) + "\n  " + m.Table.HelpView() + "\n"
	// return baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

var rootCmd = &cobra.Command{
	Use:   "ocdot",
	Short: "ocdot is dotfile managment",
	Run: func(_ *cobra.Command, _ []string) {
		homeDir, err := os.UserHomeDir()

		if err != nil {
			log.Printf("Error occured: %s", err.Error())
			os.Exit(1)
		}
		var dotfilesPath = fmt.Sprintf("%s/%s", homeDir, pathToDotfiles)
		fmt.Println(dotfilesPath)

		if err := checkIfStowExists(); err != nil {
			log.Printf("Error occured: %s", err.Error())
			os.Exit(1)
		}

		columns := []table.Column{
			{Title: "Sync", Width: 5},
			{Title: "Files", Width: 20},
			{Title: "Source", Width: 10},
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithFocused(true),
			table.WithHeight(7),
		)
		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)

		m := Model{
			Table: t,
		}
		program := tea.NewProgram(m)
		process, err := program.Run()
		if err != nil {
			log.Printf("Error occured: %s", err.Error())
			os.Exit(1)
		}
		fmt.Println(process)
	},
}

func checkIfStowExists() error {
	var isExists bool

	if _, err := os.Stat("/usr/bin/stow"); err == nil {
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
