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

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(red + "Home directory could not be retrieved. Exiting" + reset)
		return
	}
	path := filepath.Join(home, ".config", "jfri", "jfri.conf")

	if _, err = os.Stat(path); os.IsNotExist(err) {
		fmt.Println(cyan + "Config file does not exist. Creating..." + reset)
		os.MkdirAll(filepath.Dir(path), 0o755)
		file, e := os.Create(path)
		if e != nil {
			fmt.Println(red+"Error creating config file:"+reset, e)
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
		return
	}
	fmt.Println(yellow + "     ====================================" + reset)
	fmt.Println(yellow + "     =               JFRI               =" + reset)
	fmt.Println(yellow + "     ====================================" + reset)
	fmt.Println()
	// if len(extracted) == 0 {
	// 	fmt.Println("No valid entries found. Make sure the paths are in the format 'run /path/to/file'")
	// 	return
	// }

	fmt.Println(green + "Select an option:" + reset)
	fmt.Println()
	fmt.Println("[0] Edit jfri configuration file")
	for i, name := range displayNames {
		fmt.Printf("[%d] %s\n", i+1, name)
	}

	fmt.Print(yellow + "Enter the index of the option you wish to start: " + reset)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(red + "Invalid input, exiting." + reset)
		return
	}
	input = strings.TrimSpace(input)

	index, err := strconv.Atoi(input)
	if err != nil || index < 0 || index > len(extracted) {
		fmt.Println(red + "Invalid selection." + reset)
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
		fmt.Println(red + "No suitable editor found." + reset)
		return
	}

	selectedFile := extracted[index-1]
	fmt.Println(yellow+"Running:"+reset, selectedFile)

	cmd := exec.Command("sh", "-c", selectedFile) // Run as a shell command
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(red+"Error running command:"+reset, err)
	}
}
