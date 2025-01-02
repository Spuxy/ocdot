// Package tui is the main package.
package tui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spuxy/ocdot/foundation"
)

// Model is the application's main model.
type Model struct {
	Files        []string
	SyncState    string
	SelectedFile string
	Table        table.Model
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
	// rows := []table.Row{{"k", "k", "k", "k"}}
	rows, programs, err := foundation.CurrentFiles()
	if err != nil {
		panic(err)
	}
	rows = foundation.CheckProgram(rows, programs)

	m.Table.SetRows(rows)
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

// View func
func (m Model) View() string {

	// m.Table.SetHeight(len(dotfilesRows) + 1)
	m.Table.SetWidth(300)

	return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).Render(m.Table.View()) + "\n  " + m.Table.HelpView() + "\n"
	// return baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
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
