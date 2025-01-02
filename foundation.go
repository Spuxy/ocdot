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
	"golang.org/x/exp/slices"
)

// CurrentFiles func
func CurrentFiles() ([]table.Row, map[string]bool, error) {
	var dotfilesRows []table.Row
	programs := make(map[string]bool)

	rootDir := "/Users/filip.boye.kofi/dot.filesbak/"
	_, err := os.ReadDir(rootDir)
	if err != nil {
		fmt.Println(err)
		return []table.Row{}, nil, err
	}
	err = filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		dotfilesPath := strings.TrimPrefix(path, rootDir)
		// rowPath := fmt.Sprint(info.Name())
		rootProgramPosition := strings.Index(dotfilesPath, "/")
		dotfilesFile := dotfilesPath[rootProgramPosition+1:]
		program := dotfilesPath[:rootProgramPosition]
		if _, ok := programs[program]; !ok {
			programs[program] = true
		}

		dotfilesRows = append(dotfilesRows, []string{"✅", dotfilesFile, program, "github.com/spuxy/dot.file"})

		return nil
	})

	if err != nil {
		return []table.Row{}, nil, err
	}

	return dotfilesRows, programs, nil
}

// CheckProgram wtf
func CheckProgram(rows []table.Row, programs map[string]bool) []table.Row {
	// for _, value := range programs {
	//
	// }

	var finalRows []table.Row
	for _, row := range rows {
		_, err := os.Stat("/Users/filip.boye.kofi/" + row[1])
		if err != nil {
			programs[row[2]] = false
			log.Println("Error occured: ", err)
			continue
		}
	}
	for k, _ := range programs {
		var row table.Row
		if w, _ := programs[k]; w == false {
			row = table.Row{"🚫", k, "", "github.com/spuxy/dot.file"}
		} else {
			row = table.Row{"✅", k, "", "github.com/spuxy/dot.file"}
		}
		finalRows = append(finalRows, row)
	}

	slices.SortFunc(finalRows, func(i, j table.Row) int {
		return strings.Compare(i[1], j[1])
	})

	return finalRows
}

var rootCmd = &cobra.Command{
	Use:   "ocdot",
	Short: "ocdot is dotfile managment",
	Run: func(_ *cobra.Command, _ []string) {
		// homeDir, err := os.UserHomeDir()

		// if err != nil {
		// 	log.Printf("Error occured: %s", err.Error())
		// 	os.Exit(1)
		// }
		// var dotfilesPath = fmt.Sprintf("%s/%s", homeDir, pathToDotfiles)
		// fmt.Println(dotfilesPath)

		if err := checkIfStowExists(); err != nil {
			log.Printf("Error occured: %s", err.Error())
			os.Exit(1)
		}

		columns := []table.Column{
			{Title: "Sync", Width: 5},
			{Title: "Files", Width: 20},
			{Title: "Source", Width: 10},
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
		_, err := program.Run()
		if err != nil {
			log.Printf("Error occured: %s", err.Error())
			os.Exit(1)
		}
	},
}

// clearErrorAfter func
func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return ClearErrorMsg{}
	})
}

// checkIfStowExists func
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
