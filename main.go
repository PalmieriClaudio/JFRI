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
	// Start by pulling available scripts in the set folder
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Home directory could not be retrieved. exiting")
		return
	}
	path := filepath.Join(home, ".config", "jfri", "jfri.conf")

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Slice to store extracted values
	var extracted []string
	var displayNames []string

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "run ") {
			fullName := strings.TrimSpace(line[4:])
			displayName := strings.TrimSuffix(filepath.Base(fullName), filepath.Ext(fullName)) // Remove extension
			extracted = append(extracted, fullName)
			displayNames = append(displayNames, displayName)
		}
	}

	// Check for scanner errors
	if err = scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	if len(extracted) == 0 {
		fmt.Println("No valid entries found. Make sure the paths are in the format 'run /path/to/file'")
		return
	}

	// Display parsed file names
	fmt.Println("Select a file to run:")
	for i, name := range displayNames {
		fmt.Printf("[%d] %s\n", i, name)
	}

	// Get user input
	fmt.Print("Enter the index of the file to run: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Invalid input, the program will exit now.")
		return
	}
	input = strings.TrimSpace(input)

	// Convert input to an integer
	index, err := strconv.Atoi(input)
	if err != nil || index < 0 || index >= len(extracted) {
		fmt.Println("Invalid selection.")
		return
	}

	// Run the selected file
	selectedFile := extracted[index]
	if strings.HasPrefix(selectedFile, "~") {
		selectedFile = filepath.Join(home, selectedFile[1:])
	}
	// Check if the file is executable, if not make it executable
	info, err := os.Stat(selectedFile)
	if err != nil {
		fmt.Println("Error checking file:", err)
		return
	}

	// Check if the script has executable permissions
	if info.Mode()&(0o100) == 0 { // Check if the user has execute permission
		fmt.Println("Making the script executable...")
		err = os.Chmod(selectedFile, info.Mode()|0o111) // Add execute permission
		if err != nil {
			fmt.Println("Error changing file permissions:", err)
			return
		}
	}
	fmt.Println("Running:", selectedFile)

	cmd := exec.Command(selectedFile) // Assumes the file is an executable or script
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running file:", err)
	}
}
