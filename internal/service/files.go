package service

import (
	"fmt"
	"os"
)

/*
 * files.go: file to manage the reading/writing on scripts and directories.
 */

// GetDirEntries - > Gets the entries of the directory given (the path of the dir)
func GetDirEntries(name string) []os.DirEntry {
	files, err := os.ReadDir(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading the file structure: ", err)
		os.Exit(1)
	}
	return files
}

// GetWorkingDirectory Method to return the current working directory
func GetWorkingDirectory() string {
	// Get the working directory (where we are executing the commands)
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error with the path: ", err)
		os.Exit(1)
	}

	// If everything went right, we return the path of the wd
	return path
}

func ReadFile(name string) []byte {
	file, err := os.ReadFile(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading the file %s. Please report the issue.\n", name)
		os.Exit(1)
	}
	return file
}

func WriteOnFile(path string, data []byte) {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error trying to write the variable name into the file...")
		os.Exit(1)
	}
}
