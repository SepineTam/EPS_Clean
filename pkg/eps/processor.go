package eps

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/SepineTam/EPS_Clean/internal/encodingutil"
)

type FileTask struct {
	InputPath       string
	OutputPath      string
	InputEncoding   string
	OutputEncoding  string
	SafeOverwrite   bool
	CreateOutputDir bool
}

type BatchConfig struct {
	InputEncoding  string
	OutputEncoding string
}

type BatchResult struct {
	Path  string
	Error error
}

func ProcessFile(task FileTask) error {
	if task.InputPath == "" {
		return errors.New("input path cannot be empty")
	}
	if task.OutputPath == "" {
		return errors.New("output path cannot be empty")
	}

	raw, err := os.ReadFile(task.InputPath)
	if err != nil {
		return fmt.Errorf("read failed: %w", err)
	}

	decoded, err := encodingutil.Decode(raw, task.InputEncoding)
	if err != nil {
		return err
	}

	processed, err := removeLastNPhysicalLines(decoded, 3)
	if err != nil {
		return err
	}

	outBytes, err := encodingutil.Encode(processed, task.OutputEncoding)
	if err != nil {
		return err
	}

	if task.CreateOutputDir {
		if err := os.MkdirAll(filepath.Dir(task.OutputPath), 0o755); err != nil {
			return fmt.Errorf("create output directory failed: %w", err)
		}
	}

	if task.SafeOverwrite {
		return safeReplace(task.OutputPath, outBytes)
	}

	if err := os.WriteFile(task.OutputPath, outBytes, 0o644); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	return nil
}

func ProcessBatch(paths []string, cfg BatchConfig) []BatchResult {
	results := make([]BatchResult, 0, len(paths))
	for _, path := range paths {
		err := ProcessFile(FileTask{
			InputPath:      path,
			OutputPath:     path,
			InputEncoding:  cfg.InputEncoding,
			OutputEncoding: cfg.OutputEncoding,
			SafeOverwrite:  true,
		})
		results = append(results, BatchResult{Path: path, Error: err})
	}
	return results
}

func Summarize(results []BatchResult) (success int, failed int) {
	for _, r := range results {
		if r.Error != nil {
			failed++
		} else {
			success++
		}
	}
	return success, failed
}

func ScanCSVFiles(root string, recursive bool) ([]string, error) {
	entries := make([]string, 0)
	if recursive {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if isCSV(path) {
				entries = append(entries, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walk directory failed: %w", err)
		}
		sort.Strings(entries)
		return entries, nil
	}

	dirEntries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read directory failed: %w", err)
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(root, entry.Name())
		if isCSV(path) {
			entries = append(entries, path)
		}
	}
	sort.Strings(entries)
	return entries, nil
}

func removeLastNPhysicalLines(text string, n int) (string, error) {
	lines := splitPhysicalLines(text)
	if len(lines) <= n {
		return "", fmt.Errorf("file has %d physical lines; at least %d required", len(lines), n+1)
	}
	return strings.Join(lines[:len(lines)-n], ""), nil
}

func splitPhysicalLines(text string) []string {
	if text == "" {
		return []string{""}
	}
	lines := make([]string, 0)
	start := 0
	for i := 0; i < len(text); i++ {
		if text[i] == '\n' {
			lines = append(lines, text[start:i+1])
			start = i + 1
		}
	}
	if start < len(text) {
		lines = append(lines, text[start:])
	}
	return lines
}

func isCSV(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".csv")
}

func safeReplace(target string, data []byte) error {
	dir := filepath.Dir(target)
	base := filepath.Base(target)
	tmp, err := os.CreateTemp(dir, base+".tmp-*")
	if err != nil {
		return fmt.Errorf("create temp file failed: %w", err)
	}
	tmpPath := tmp.Name()

	writeErr := func() error {
		if _, err := tmp.Write(data); err != nil {
			return fmt.Errorf("write temp file failed: %w", err)
		}
		if err := tmp.Sync(); err != nil {
			return fmt.Errorf("sync temp file failed: %w", err)
		}
		if err := tmp.Close(); err != nil {
			return fmt.Errorf("close temp file failed: %w", err)
		}
		return nil
	}()
	if writeErr != nil {
		_ = os.Remove(tmpPath)
		return writeErr
	}

	if err := os.Rename(tmpPath, target); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("replace original file failed: %w", err)
	}
	return nil
}
