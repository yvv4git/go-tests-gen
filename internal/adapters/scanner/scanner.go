package scanner

import (
	"bufio"
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yvv4git/go-tests-gen/internal/ports"
)

type Coverage struct {
	path           string
	uncoveredFuncs []ports.UncoveredFunc
}

func NewCoverage(path string) *Coverage {
	return &Coverage{
		path:           path,
		uncoveredFuncs: make([]ports.UncoveredFunc, 0),
	}
}

func (c *Coverage) Scan(ctx context.Context) error {
	coverProfile := c.path + "/cover.out"

	cmd := exec.CommandContext(ctx, "go", "test", "-coverprofile="+coverProfile, "-covermode=atomic", "-cover", "./...")
	cmd.Dir = c.path

	// Ignore test failures, we only need coverage data
	cmd.Run()

	if _, err := os.Stat(coverProfile); os.IsNotExist(err) {
		return nil
	}

	coveredRanges, err := c.parseCoverProfile(coverProfile)
	if err != nil {
		os.Remove(coverProfile)
		return err
	}

	c.uncoveredFuncs, err = c.findUncoveredFunctions(coveredRanges)
	if err != nil {
		os.Remove(coverProfile)
		return err
	}

	os.Remove(coverProfile)

	return nil
}

func (c *Coverage) parseCoverProfile(coverProfile string) (map[string][][2]int, error) {
	coveredRanges := make(map[string][][2]int)

	file, err := os.Open(coverProfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "mode:") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		location := parts[0]
		countStr := parts[2]
		count, err := strconv.Atoi(countStr)
		if err != nil {
			continue
		}

		if count == 0 {
			pos := strings.Split(location, ":")
			if len(pos) != 2 {
				continue
			}

			// Format: startLine.startCol,endLine.endCol
			// Split by comma to get start and end positions
			rangeParts := strings.Split(pos[1], ",")
			if len(rangeParts) != 2 {
				continue
			}

			// Each position is in format line.column
			startPos := strings.Split(rangeParts[0], ".")
			endPos := strings.Split(rangeParts[1], ".")
			if len(startPos) != 2 || len(endPos) != 2 {
				continue
			}

			startLine, _ := strconv.Atoi(startPos[0])
			endLine, err := strconv.Atoi(endPos[0])
			if err != nil {
				continue
			}

			filePath := c.path + "/" + pos[0]
			// Add all lines in the range [startLine, endLine]
			for line := startLine; line <= endLine; line++ {
				coveredRanges[filePath] = append(coveredRanges[filePath], [2]int{line, line})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return coveredRanges, nil
}

func (c *Coverage) findUncoveredFunctions(coveredRanges map[string][][2]int) ([]ports.UncoveredFunc, error) {
	var uncovered []ports.UncoveredFunc

	err := filepath.Walk(c.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil
		}

		for _, decl := range node.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			funcLine := fset.Position(fn.Pos()).Line
			funcEndLine := fset.Position(fn.End()).Line

			isCovered := false
			for _, rangePair := range coveredRanges[path] {
				for i := funcLine; i <= funcEndLine; i++ {
					if i >= rangePair[0] && i <= rangePair[1] {
						isCovered = true
						break
					}
				}
				if isCovered {
					break
				}
			}

			if !isCovered {
				relPath, _ := filepath.Rel(c.path, path)
				pkg := c.getPackageName(path)
				fnCode := c.extractFuncCode(path, funcLine, funcEndLine)
				uncovered = append(uncovered, ports.UncoveredFunc{
					Package: pkg,
					File:    relPath,
					Name:    fn.Name.Name,
					Line:    funcLine,
					FnCode:  fnCode,
				})
			}
		}

		return nil
	})

	return uncovered, err
}

func (c *Coverage) getPackageName(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if after, ok := strings.CutPrefix(line, "package "); ok {
			return after
		}
	}

	return ""
}

func (c *Coverage) extractFuncCode(filePath string, startLine, endLine int) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for lineNum := 1; scanner.Scan(); lineNum++ {
		if lineNum >= startLine && lineNum <= endLine {
			lines = append(lines, scanner.Text())
		}
	}

	return strings.Join(lines, "\n")
}

func (c *Coverage) GetUncoveredFunctions() []ports.UncoveredFunc {
	return c.uncoveredFuncs
}
