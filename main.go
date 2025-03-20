package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Home directory could not be retrieved. Exiting")
		return
	}
	path := filepath.Join(home, ".config", "jfri", "jfri.conf")

	if _, err = os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Config file does not exist. Creating...")
		os.MkdirAll(filepath.Dir(path), 0o755)
		file, e := os.Create(path)
		if e != nil {
			fmt.Println("Error creating config file:", e)
			return
		}
		file.Close()
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var extracted []string
	var displayNames []string
	var currentName string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "name ") {
			currentName = strings.TrimSpace(line[5:])
		} else if strings.HasPrefix(line, "run ") {
			fullCmd := strings.TrimSpace(line[4:])
			if currentName == "" {
				currentName = strings.TrimSuffix(filepath.Base(fullCmd), filepath.Ext(fullCmd))
			}
			extracted = append(extracted, fullCmd)
			displayNames = append(displayNames, currentName)
			currentName = ""
		}
	}

	if err = scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	// if len(extracted) == 0 {
	// 	fmt.Println("No valid entries found. Make sure the paths are in the format 'run /path/to/file'")
	// 	return
	// }

	fmt.Println("Select an option:")
	fmt.Println("[0] Open/Edit jfri configuration file")
	for i, name := range displayNames {
		fmt.Printf("[%d] %s\n", i+1, name)
	}

	fmt.Print("Enter the index of the option: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Invalid input, exiting.")
		return
	}
	input = strings.TrimSpace(input)

	index, err := strconv.Atoi(input)
	if err != nil || index < 0 || index > len(extracted) {
		fmt.Println("Invalid selection.")
		return
	}

	if index == 0 {
		editors := []string{"nvim", "nano"}
		for _, editor := range editors {
			cmd := exec.Command("which", editor)
			if err = cmd.Run(); err == nil {
				cmd = exec.Command(editor, path)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err = cmd.Run(); err == nil {
					return
				}
			}
		}
		fmt.Println("No suitable editor found.")
		return
	}

	selectedFile := extracted[index-1]
	fmt.Println("Running:", selectedFile)

	cmd := exec.Command("sh", "-c", selectedFile) // Run as a shell command
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running command:", err)
	}
}
