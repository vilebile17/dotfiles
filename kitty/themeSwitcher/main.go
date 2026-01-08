package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

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

	allFiles, err := os.ReadDir(homeDir + "/.config/kitty")
	if err != nil {
		log.Fatal(err)
	}

	configFileNames := []string{}
	for _, file := range allFiles {
		hiddenFile := strings.HasPrefix(file.Name(), ".")
		if file.Name() != "kitty.conf" && !hiddenFile && !file.IsDir() {
			nameWithoutSuffix := strings.Split(file.Name(), ".")[0]
			configFileNames = append(configFileNames, nameWithoutSuffix)
		}
	}

	lines, err := OpenAndReadConfig(homeDir)
	if err != nil {
		log.Fatal(err)
	}
	filePathParts := strings.Split(lines[len(lines)-1], "/")
	currentFile := filePathParts[len(filePathParts)-1]

	selected := 0
	for i, file := range configFileNames {
		if file+".conf" == currentFile {
			selected = i
		}
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

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.themes)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			m.selected = m.cursor

			// This part Adds the file to the kitty.conf file and execute 'kitty @ load-config'
			m.AddToConfig(m.themes[m.selected] + ".conf")
			cmd := exec.Command("kitty", "@", "load-config")
			go cmd.Run()
		}
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func OpenAndReadConfig(homeDir string) ([]string, error) {
	file, err := os.Open(homeDir + "/.config/kitty/kitty.conf")
	if err != nil {
		return []string{}, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	return lines, nil
}

func (m model) AddToConfig(filename string) {
	lines, err := OpenAndReadConfig(m.homeDir)
	if err != nil {
		log.Fatal(err)
	}

	lines[len(lines)-1] = "include ~/.config/kitty/" + filename
	var buffer bytes.Buffer
	for _, line := range lines {
		buffer.WriteString(line + "\n")
	}

	err = os.WriteFile(m.homeDir+"/.config/kitty/kitty.conf", buffer.Bytes(), 0o644)
	if err != nil {
		log.Fatal(err)
	}
}

func (m model) View() string {
	// The header
	s := "What theme would you like?\n\n"

	// Iterate over our choices
	for i, theme := range m.themes {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
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
