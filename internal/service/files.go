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
		return nil
	}
	return files
}

// GetWorkingDirectory Method to return the current working directory
func GetWorkingDirectory() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting working directory: ", err)
		return "."
	}
	return path
}

func ReadFile(name string) []byte {
	file, err := os.ReadFile(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading the file %s: %v\n", name, err)
		return nil
	}
	return file
}

func WriteOnFile(path string, data []byte) {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file %s: %v\n", path, err)
	}
}

