package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

type model struct {
	themes   []string
	cursor   int
	selected int
	homeDir  string
}

func initialModel() model {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	configFileNames, err := getConfigFiles(homeDir)
	if err != nil {
		log.Fatal(err)
	}

	selected, err := findCurrentFileIndex(configFileNames, homeDir)
	if err != nil {
		log.Fatal(err)
	}

	return model{
		themes:   configFileNames,
		cursor:   selected,
		selected: selected,
		homeDir:  homeDir,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.themes) - 1
			}
		case "down", "j":
			if m.cursor < len(m.themes)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "enter", " ":
			m.selected = m.cursor

			err := m.AddToConfig(m.themes[m.selected] + ".conf")
			if err != nil {
				log.Fatal(err)
			}

			cmd := exec.Command("kitty", "@", "load-config")
			go cmd.Run()
		}
	}
	return m, nil
}

func (m model) View() string {
	// The header
	s := "What theme would you like?\n\n"

	for i, theme := range m.themes {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		checked := " " // not selected
		if i == m.selected {
			checked = "x"
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, theme)
	}

	// The footer
	s += color.RedString("\nRed ")
	s += color.YellowString("Yellow ")
	s += color.GreenString("Green ")
	s += color.CyanString("Cyan ")
	s += color.BlueString("Blue ")
	s += color.MagentaString("Magenta ")
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
