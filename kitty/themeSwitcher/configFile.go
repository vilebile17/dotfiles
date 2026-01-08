package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func getConfigFiles(homeDir string) ([]string, error) {
	allFiles, err := os.ReadDir(homeDir + "/.config/kitty")
	if err != nil {
		return []string{}, err
	}

	configFileNames := []string{}
	for _, file := range allFiles {
		hiddenFile := strings.HasPrefix(file.Name(), ".")
		if file.Name() != "kitty.conf" && !hiddenFile && !file.IsDir() {
			nameWithoutSuffix := strings.Split(file.Name(), ".")[0]
			configFileNames = append(configFileNames, nameWithoutSuffix)
		}
	}
	return configFileNames, err
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

func (m model) AddToConfig(filename string) error {
	lines, err := OpenAndReadConfig(m.homeDir)
	if err != nil {
		return err
	}

	lines[len(lines)-1] = "include ~/.config/kitty/" + filename
	var buffer bytes.Buffer
	for _, line := range lines {
		buffer.WriteString(line + "\n")
	}

	err = os.WriteFile(m.homeDir+"/.config/kitty/kitty.conf", buffer.Bytes(), 0o644)
	if err != nil {
		return err
	}
	return nil
}

func findCurrentFileIndex(files []string, homeDir string) (int, error) {
	lines, err := OpenAndReadConfig(homeDir)
	if err != nil {
		return 0, err
	}
	filePathParts := strings.Split(lines[len(lines)-1], "/")
	currentFile := filePathParts[len(filePathParts)-1]

	for i := range files {
		if files[i]+".conf" == currentFile {
			return i, nil
		}
	}
	return 0, fmt.Errorf("couldn't find file: %v in the theme directory", currentFile)
}
