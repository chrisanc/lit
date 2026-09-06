package service

import (
	"CLI_App/internal/domain"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// ScanService defines a service (use case) for the file scanning
type ScanService struct {
	wg                 sync.WaitGroup
	mu                 sync.Mutex
	analyzer           domain.Analyzer
	dangerousFunctions map[string][]*domain.FunctionData
	languagesMap       map[string]int
}

func NewScannerService(analyzer domain.Analyzer) ScanService {
	return ScanService{
		analyzer:           analyzer,
		dangerousFunctions: make(map[string][]*domain.FunctionData),
		languagesMap:       make(map[string]int),
	}
}

// ScanFiles starts the scanning process. Entry point
func (s *ScanService) ScanFiles() {
	s.traverseFiles(s.scanFile, domain.ScanValidScriptPattern)
}

// ExecuteLOC starts the scanning process for the loc data. Entry point.
func (s *ScanService) ExecuteLOC() {
	s.traverseFiles(s.loc, domain.LocValidScriptPattern)
}

// FixFile fixes the name of certain variables
func (s *ScanService) FixFile() {
	s.traverseFiles(s.fixFile, domain.ScanValidScriptPattern)
}

// Internal functions to analyze code
func (s *ScanService) scanFile(filename string, code *[]string) {
	defer s.wg.Done()
	if code == nil {
		return
	}
	// Execute the file analyzer algorithm (infrastructure because it's using frameworks)
	functions := s.analyzer.AnalyzeFile(filename, code)
	if functions != nil && len(functions) > 0 {
		s.mu.Lock()
		s.dangerousFunctions[filename] = functions
		s.mu.Unlock()
	}
}

func (s *ScanService) loc(filename string, code *[]string) {
	defer s.wg.Done()
	// Sum up the stored value with the total lines found in that script.
	s.mu.Lock()
	s.languagesMap[filepath.Ext(filename)[1:]] += len(*code)
	s.mu.Unlock()
}

func (s *ScanService) fixFile(filename string, code *[]string) {
	defer s.wg.Done()
	s.languagesMap[filename] += s.analyzer.FixFile(filename, code)
	WriteOnFile(filename, []byte(strings.Join(*code, "\n")))
}

// Navigate through the file system with a DFS algorithm.
func (s *ScanService) traverseFiles(fileFunction func(filename string, code *[]string), validScriptPattern string) {
	stack := []domain.Directory{{"", GetDirEntries(GetWorkingDirectory())}}
	for len(stack) > 0 {
		// Extract the last element from the stack
		files := stack[len(stack)-1]
		// Remove the last element from the stack (the files we just iterated).
		stack = stack[:len(stack)-1]
		for _, v := range files.Content {
			// Check out if the current position contains a file or a directory
			if v.IsDir() {
				// If we should ignore a directory based on our regex, we do.
				if r, _ := regexp.Match(domain.NotValidDirPattern, []byte(v.Name())); r {
					continue
				}
				fmt.Println("Reading", files.DirName+v.Name()+"/")
				dir := GetDirEntries(files.DirName + v.Name() + "/")
				stack = append(stack, domain.Directory{DirName: files.DirName + v.Name() + "/", Content: dir})
			} else {
				// Check if the current file is a programming language script
				if r, _ := regexp.Match(validScriptPattern, []byte(v.Name())); !r {
					continue
				}
				// Create the file path
				path := files.DirName + v.Name()

				file := strings.Split(string(ReadFile(path)), "\n")

				s.wg.Add(1)
				go fileFunction(path, &file)
			}
		}
	}
	s.wg.Wait()
}

// PrintLOCResults prints out the results of the loc execution
func (s *ScanService) PrintLOCResults() {
	var totalLines float64
	fmt.Println()
	fmt.Println("Results (language - > total lines of code)")

	for _, v := range s.languagesMap {
		totalLines += float64(v)
	}

	for key, value := range s.languagesMap {
		fmt.Printf("%s - > %d (%.1f%%)\n", key, value, (float64(value)*100)/totalLines)
	}
}

func (s *ScanService) PrintScanningResults() {
	totalFunctions := 0
	for _, v := range s.dangerousFunctions {
		totalFunctions += len(v)
	}

	fmt.Println()
	fmt.Printf("Found %d possible improvements\n", totalFunctions)
	for key, value := range s.dangerousFunctions {
		fmt.Printf("- %s:\n", key)
		for _, item := range value {
			fmt.Printf(" * %s (at %d:%d)\n", item.Name, item.StartPosition.Row, item.StartPosition.Column)
			fmt.Printf("   Parameters: %d\n   Total lines of code: %d\n", item.TotalParams, item.Size)
			fmt.Printf("   Found %d variables/methods/models with the wrong naming convention.\n", item.InvalidNames)
			fmt.Println(item.Feedback)
		}
	}
}

func (s *ScanService) PrintFixResults() {
	for key, value := range s.languagesMap {
		fmt.Printf("%s - > %d names modified.\n", key, value)
	}
}
