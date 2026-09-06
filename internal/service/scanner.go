package service

import (
	"CLI_App/internal/adapters/exporter"
	"CLI_App/internal/domain"
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// ScanService defines a service (use case) for the file scanning
type ScanService struct {
	mu                 sync.Mutex
	analyzer           domain.Analyzer
	dangerousFunctions map[string][]*domain.FunctionData
	languagesMap       map[string]int
	diffsMap           map[string]string
	workersCount       int
}

func NewScannerService(analyzer domain.Analyzer) ScanService {
	workers := runtime.NumCPU() * 2
	if workers < 4 {
		workers = 4
	}
	return ScanService{
		analyzer:           analyzer,
		dangerousFunctions: make(map[string][]*domain.FunctionData),
		languagesMap:       make(map[string]int),
		diffsMap:           make(map[string]string),
		workersCount:       workers,
	}
}

// SetWorkers allows configuring the worker pool concurrency size
func (s *ScanService) SetWorkers(count int) {
	if count > 0 {
		s.workersCount = count
	}
}

// ScanFiles starts the scanning process with context support. Entry point.
func (s *ScanService) ScanFiles(ctx context.Context) {
	s.traverseFiles(ctx, s.scanFile, domain.ScanValidScriptRegexp)
}

// ExecuteLOC starts the scanning process for the loc data with context support. Entry point.
func (s *ScanService) ExecuteLOC(ctx context.Context) {
	s.traverseFiles(ctx, s.loc, domain.LocValidScriptRegexp)
}

// FixFile fixes the name of certain variables with context and dryRun support.
func (s *ScanService) FixFile(ctx context.Context, dryRun bool) {
	if dryRun {
		s.traverseFiles(ctx, s.fixFileDryRun, domain.ScanValidScriptRegexp)
	} else {
		s.traverseFiles(ctx, s.fixFile, domain.ScanValidScriptRegexp)
	}
}

// Internal functions to analyze code
func (s *ScanService) scanFile(filename string, code *[]string) {
	if code == nil {
		return
	}
	functions := s.analyzer.AnalyzeFile(filename, code)
	if len(functions) > 0 {
		s.mu.Lock()
		s.dangerousFunctions[filename] = functions
		s.mu.Unlock()
	}
}

func (s *ScanService) loc(filename string, code *[]string) {
	if code == nil {
		return
	}
	ext := filepath.Ext(filename)
	if len(ext) > 1 {
		s.mu.Lock()
		s.languagesMap[ext[1:]] += len(*code)
		s.mu.Unlock()
	}
}

func (s *ScanService) fixFile(filename string, code *[]string) {
	if code == nil {
		return
	}
	modified := s.analyzer.FixFile(filename, code)
	if modified > 0 {
		s.mu.Lock()
		s.languagesMap[filename] += modified
		s.mu.Unlock()
		WriteOnFile(filename, []byte(strings.Join(*code, "\n")))
	}
}

func (s *ScanService) fixFileDryRun(filename string, code *[]string) {
	if code == nil {
		return
	}
	// Copy original code lines
	original := make([]string, len(*code))
	copy(original, *code)

	modified := s.analyzer.FixFile(filename, code)
	if modified > 0 {
		diff := GenerateUnifiedDiff(filename, original, *code)
		s.mu.Lock()
		s.languagesMap[filename] += modified
		s.diffsMap[filename] = diff
		s.mu.Unlock()
	}
}

// Navigate through the file system with a bounded worker pool and context cancellation.
func (s *ScanService) traverseFiles(ctx context.Context, fileFunction func(filename string, code *[]string), validScriptRegexp *regexp.Regexp) {
	jobs := make(chan string, 100)
	var workerWg sync.WaitGroup

	// Launch bounded worker pool
	for i := 0; i < s.workersCount; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case path, ok := <-jobs:
					if !ok {
						return
					}
					rawBytes := ReadFile(path)
					if rawBytes == nil {
						continue
					}
					lines := strings.Split(string(rawBytes), "\n")
					fileFunction(path, &lines)
				}
			}
		}()
	}

	// DFS directory traversal pushing file paths to worker queue
	stack := []domain.Directory{{"", GetDirEntries(GetWorkingDirectory())}}
Loop:
	for len(stack) > 0 {
		select {
		case <-ctx.Done():
			break Loop
		default:
		}

		files := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, v := range files.Content {
			select {
			case <-ctx.Done():
				break Loop
			default:
			}

			if v.IsDir() {
				if domain.NotValidDirRegexp.MatchString(v.Name()) {
					continue
				}
				dirPath := files.DirName + v.Name() + "/"
				dir := GetDirEntries(dirPath)
				if dir != nil {
					stack = append(stack, domain.Directory{DirName: dirPath, Content: dir})
				}
			} else {
				if !validScriptRegexp.MatchString(v.Name()) {
					continue
				}
				path := files.DirName + v.Name()
				select {
				case jobs <- path:
				case <-ctx.Done():
					break Loop
				}
			}
		}
	}

	close(jobs)
	workerWg.Wait()
}

// PrintLOCResults prints out the results of the loc execution
func (s *ScanService) PrintLOCResults() {
	var totalLines float64
	fmt.Println()
	fmt.Println("Results (language -> total lines of code)")

	for _, v := range s.languagesMap {
		totalLines += float64(v)
	}

	if totalLines == 0 {
		fmt.Println("No matching files found.")
		return
	}

	for key, value := range s.languagesMap {
		fmt.Printf("%s -> %d (%.1f%%)\n", key, value, (float64(value)*100)/totalLines)
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

func (s *ScanService) GetDangerousFunctions() map[string][]*domain.FunctionData {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.dangerousFunctions
}

func (s *ScanService) ExportResults(format, outputPath string) error {
	exp, err := exporter.GetExporter(format)
	if err != nil {
		return err
	}

	data, err := exp.Export(s.GetDangerousFunctions())
	if err != nil {
		return fmt.Errorf("error generating %s report: %w", format, err)
	}

	if outputPath != "" {
		WriteOnFile(outputPath, data)
		fmt.Printf("Report exported successfully to %s (%s format)\n", outputPath, format)
	} else {
		fmt.Println(string(data))
	}

	return nil
}

func (s *ScanService) PrintFixResults(dryRun bool) {
	if dryRun {
		fmt.Println()
		fmt.Println("--- DRY RUN: Unified Git Diff Preview ---")
		totalFiles := 0
		totalMods := 0
		for file, diff := range s.diffsMap {
			totalFiles++
			totalMods += s.languagesMap[file]
			fmt.Println(diff)
		}
		if totalFiles == 0 {
			fmt.Println("No variable naming fixes required.")
		} else {
			fmt.Printf("[Dry Run Complete] Previewed %d variable naming changes across %d file(s). No files were modified on disk.\n", totalMods, totalFiles)
		}
	} else {
		for key, value := range s.languagesMap {
			fmt.Printf("%s -> %d names modified.\n", key, value)
		}
	}
}
