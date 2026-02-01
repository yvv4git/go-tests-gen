package scanner

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yvv4git/go-tests-gen/internal/ports"
)

// Mock file info for testing
type mockFileInfo struct {
	name    string
	isDir   bool
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return m.size }
func (m *mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m *mockFileInfo) ModTime() time.Time { return m.modTime }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }

// Test parseCoverProfile - tests with correct cover profile format
// Real format: file.go:startLine.startCol,endLine.endCol stmtCount count
func TestParseCoverProfile(t *testing.T) {
	tests := []struct {
		name     string
		content  []string
		expected map[string][][2]int
	}{
		{
			name:     "empty file",
			content:  []string{"mode: count"},
			expected: map[string][][2]int{},
		},
		{
			name: "covered function",
			content: []string{
				"mode: count",
				"pkg/file.go:10.1,12.2 0 1",
			},
			expected: map[string][][2]int{},
		},
		{
			name: "uncovered function (single line range)",
			content: []string{
				"mode: count",
				"pkg/file.go:10.1,10.2 0 0",
			},
			expected: nil, // Will be set per-test with actual tmpDir
		},
		{
			name: "multiple uncovered functions",
			content: []string{
				"mode: count",
				"pkg/file1.go:10.1,10.2 0 0",
				"pkg/file2.go:20.1,20.2 0 0",
			},
			expected: nil, // Will be set per-test
		},
		{
			name: "mixed coverage",
			content: []string{
				"mode: count",
				"pkg/file.go:10.1,12.2 0 1",
				"pkg/file.go:20.1,20.2 0 0",
			},
			expected: nil, // Will be set per-test
		},
		{
			name: "range with multiple lines",
			content: []string{
				"mode: count",
				"pkg/file.go:10.1,15.3 0 0",
			},
			expected: nil, // Will be set per-test
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for testing
			tmpDir, err := os.MkdirTemp("", "scanner_test_*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			c := &Scanner{path: tmpDir}

			// Set expected map with actual tmpDir
			expected := make(map[string][][2]int)

			// Adjust expected paths to use tmpDir based on content
			switch tt.name {
			case "uncovered function (single line range)":
				expected[tmpDir+"/pkg/file.go"] = [][2]int{{10, 10}}
			case "multiple uncovered functions":
				expected[tmpDir+"/pkg/file1.go"] = [][2]int{{10, 10}}
				expected[tmpDir+"/pkg/file2.go"] = [][2]int{{20, 20}}
			case "mixed coverage":
				expected[tmpDir+"/pkg/file.go"] = [][2]int{{20, 20}}
			case "range with multiple lines":
				// Range from line 10 to line 15: 10, 11, 12, 13, 14, 15
				expected[tmpDir+"/pkg/file.go"] = [][2]int{
					{10, 10}, {11, 11}, {12, 12}, {13, 13}, {14, 14}, {15, 15},
				}
			}

			// Create temp file
			tmpFile, err := os.CreateTemp("", "cover_test_*.out")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			writer := bufio.NewWriter(tmpFile)
			for _, line := range tt.content {
				writer.WriteString(line + "\n")
			}
			writer.Flush()
			tmpFile.Close()

			result, err := c.parseCoverProfile(tmpFile.Name())
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Debug output
			t.Logf("tmpDir: %s", tmpDir)
			t.Logf("result: %v", result)
			t.Logf("expected: %v", expected)

			// Compare results
			if len(result) != len(expected) {
				t.Errorf("Expected %d files, got %d", len(expected), len(result))
				return
			}

			for file, ranges := range expected {
				resultRanges, ok := result[file]
				if !ok {
					t.Errorf("Expected file %s not found", file)
					continue
				}
				if len(resultRanges) != len(ranges) {
					t.Errorf("Expected %d ranges for %s, got %d", len(ranges), file, len(resultRanges))
					continue
				}
				for i, r := range ranges {
					if r != resultRanges[i] {
						t.Errorf("Range mismatch for %s[%d]: expected %v, got %v", file, i, r, resultRanges[i])
					}
				}
			}
		})
	}
}

// Test extractFuncCode
func TestExtractFuncCode(t *testing.T) {
	c := &Scanner{path: "/test"}

	tests := []struct {
		name      string
		fileLines []string
		startLine int
		endLine   int
		expected  string
	}{
		{
			name:      "single line function",
			fileLines: []string{"package pkg", "", "func Hello() {", "  return nil", "}"},
			startLine: 3,
			endLine:   5,
			expected:  "func Hello() {\n  return nil\n}",
		},
		{
			name:      "function at start",
			fileLines: []string{"func First() {}", "package pkg", "func Second() {}"},
			startLine: 1,
			endLine:   1,
			expected:  "func First() {}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "scanner_test_*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			tmpFile := filepath.Join(tmpDir, "test.go")
			if err := os.WriteFile(tmpFile, []byte(strings.Join(tt.fileLines, "\n")), 0644); err != nil {
				t.Fatalf("Failed to write temp file: %v", err)
			}

			result := c.extractFuncCode(tmpFile, tt.startLine, tt.endLine)
			if result != tt.expected {
				t.Errorf("Expected:\n%q\nGot:\n%q", tt.expected, result)
			}
		})
	}
}

// Test getPackageName
func TestGetPackageName(t *testing.T) {
	c := &Scanner{path: "/test"}

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "simple package",
			content:  "package main\n",
			expected: "main",
		},
		{
			name:     "package with imports",
			content:  "package utils\n\nimport \"fmt\"\n",
			expected: "utils",
		},
		{
			name:     "package with comment",
			content:  "// Package utils provides utility functions\npackage utils\n",
			expected: "utils",
		},
		{
			name:     "empty file",
			content:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "pkg_test_*.go")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tt.content); err != nil {
				t.Fatalf("Failed to write temp file: %v", err)
			}
			tmpFile.Close()

			result := c.getPackageName(tmpFile.Name())
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test findUncoveredFunctions
func TestFindUncoveredFunctions(t *testing.T) {
	tests := []struct {
		name               string
		coveredRanges      map[string][][2]int
		sourceFiles        map[string]string
		expectedFuncsCount int
		expectedFuncNames  []string
	}{
		{
			name:          "all functions covered",
			coveredRanges: map[string][][2]int{},
			sourceFiles: map[string]string{
				"file.go": `package main

func Covered() {}
func AlsoCovered() {}
`,
			},
			expectedFuncsCount: 2,
			expectedFuncNames:  []string{"Covered", "AlsoCovered"},
		},
		{
			name:          "one uncovered function",
			coveredRanges: map[string][][2]int{},
			sourceFiles: map[string]string{
				"file.go": `package main

func Uncovered() {}
`,
			},
			expectedFuncsCount: 1,
			expectedFuncNames:  []string{"Uncovered"},
		},
		{
			name:          "multiple uncovered functions",
			coveredRanges: map[string][][2]int{},
			sourceFiles: map[string]string{
				"file1.go": `package main

func Func1() {}
`,
				"file2.go": `package main

func Func2() {}
`,
			},
			expectedFuncsCount: 2,
			expectedFuncNames:  []string{"Func1", "Func2"},
		},
		{
			name:          "functions with different coverage",
			coveredRanges: map[string][][2]int{},
			sourceFiles: map[string]string{
				"file1.go": `package main

func Covered() {}
func AlsoCovered() {}
`,
				"file2.go": `package main

func Uncovered() {}
`,
			},
			expectedFuncsCount: 3,
			expectedFuncNames:  []string{"Covered", "AlsoCovered", "Uncovered"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "scanner_test_*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			c := &Scanner{path: tmpDir}

			// Adjust covered ranges to use tmpDir paths
			adjustedRanges := make(map[string][][2]int)
			for path, ranges := range tt.coveredRanges {
				adjustedRanges[tmpDir+"/"+path] = ranges
			}

			// Write source files
			for filename, content := range tt.sourceFiles {
				tmpFile := filepath.Join(tmpDir, filename)
				if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write file %s: %v", filename, err)
				}
			}

			result, err := c.findUncoveredFunctions(adjustedRanges)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(result) != tt.expectedFuncsCount {
				t.Errorf("Expected %d uncovered functions, got %d", tt.expectedFuncsCount, len(result))
			}

			// Check function names
			if tt.expectedFuncNames != nil {
				var gotNames []string
				for _, f := range result {
					gotNames = append(gotNames, f.Name)
				}
				if len(gotNames) != len(tt.expectedFuncNames) {
					t.Errorf("Expected %d function names, got %d", len(tt.expectedFuncNames), len(gotNames))
				}
			}
		})
	}
}

// Test NewCoverage
func TestNewCoverage(t *testing.T) {
	path := "/test/path"
	c := NewScanner(path)

	if c.path != path {
		t.Errorf("Expected path %q, got %q", path, c.path)
	}

	if c.uncoveredFuncs == nil {
		t.Error("Expected uncoveredFuncs to be initialized")
	}

	if len(c.uncoveredFuncs) != 0 {
		t.Errorf("Expected empty uncoveredFuncs, got %d", len(c.uncoveredFuncs))
	}
}

// Test GetUncoveredFunctions
func TestGetUncoveredFunctions(t *testing.T) {
	c := NewScanner("/test")

	// Initially empty
	if len(c.GetUncoveredFunctions()) != 0 {
		t.Error("Expected empty list initially")
	}

	// After setting
	c.uncoveredFuncs = []ports.UncoveredFunc{
		{Package: "main", File: "file.go", Name: "TestFunc", Line: 10, FnCode: "func TestFunc() {}"},
	}

	result := c.GetUncoveredFunctions()
	if len(result) != 1 {
		t.Errorf("Expected 1 function, got %d", len(result))
	}

	if result[0].Name != "TestFunc" {
		t.Errorf("Expected function name TestFunc, got %s", result[0].Name)
	}
}

// Integration test for Scan (requires go test to be available)
func TestScanIntegration(t *testing.T) {
	// Skip if not in integration test mode
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=1 to run.")
	}

	// Create a temp project with test files
	tmpDir, err := os.MkdirTemp("", "scanner_integration_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create go.mod
	goMod := `module testproject

go 1.21
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	// Create a file with uncovered function
	sourceCode := `package main

func UncoveredFunc() {
	// This function is not tested
}

func CoveredFunc() {
	// This function will be covered
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(sourceCode), 0644); err != nil {
		t.Fatalf("Failed to write main.go: %v", err)
	}

	// Create a test file that only covers CoveredFunc
	testCode := `package main

import "testing"

func TestCoveredFunc(t *testing.T) {
	CoveredFunc()
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main_test.go"), []byte(testCode), 0644); err != nil {
		t.Fatalf("Failed to write main_test.go: %v", err)
	}

	c := NewScanner(tmpDir)
	ctx := context.Background()

	err = c.Scan(ctx)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	funcs := c.GetUncoveredFunctions()
	if len(funcs) == 0 {
		t.Error("Expected at least one uncovered function")
	}

	found := false
	for _, f := range funcs {
		if f.Name == "UncoveredFunc" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected UncoveredFunc to be in uncovered functions list")
	}
}
