package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SepineTam/EPS_Clean/pkg/eps"
)

type options struct {
	InputEncoding  string
	OutputEncoding string
	FromDir        string
	Recursive      bool
	Help           bool
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printRootHelp()
		return nil
	}

	if args[0] == "-h" || args[0] == "--help" {
		printRootHelp()
		return nil
	}

	if args[0] == "batch" {
		return runBatch(args[1:])
	}

	return runSingle(args)
}

func runSingle(args []string) error {
	opts, positional, err := parseOptions(args, false)
	if err != nil {
		return err
	}
	if opts.Help {
		printSingleHelp()
		return nil
	}

	if len(positional) == 0 || len(positional) > 2 {
		return errors.New("single-file mode requires 1 or 2 file arguments")
	}

	input := positional[0]
	output := input
	if len(positional) == 2 {
		output = positional[1]
	}

	if err := eps.ProcessFile(eps.FileTask{
		InputPath:       input,
		OutputPath:      output,
		InputEncoding:   opts.InputEncoding,
		OutputEncoding:  opts.OutputEncoding,
		SafeOverwrite:   input == output,
		CreateOutputDir: false,
	}); err != nil {
		return err
	}

	fmt.Printf("Processed %s -> %s\n", input, output)
	return nil
}

func runBatch(args []string) error {
	opts, positional, err := parseOptions(args, true)
	if err != nil {
		return err
	}
	if opts.Help {
		printBatchHelp()
		return nil
	}

	inputs := make([]string, 0, len(positional))
	inputs = append(inputs, positional...)

	if opts.FromDir != "" {
		dirFiles, err := eps.ScanCSVFiles(opts.FromDir, opts.Recursive)
		if err != nil {
			return err
		}
		inputs = append(inputs, dirFiles...)
	}

	if len(inputs) == 0 {
		return errors.New("batch mode requires file arguments and/or --from-dir")
	}

	results := eps.ProcessBatch(inputs, eps.BatchConfig{
		InputEncoding:  opts.InputEncoding,
		OutputEncoding: opts.OutputEncoding,
	})

	for _, result := range results {
		if result.Error != nil {
			fmt.Fprintf(os.Stderr, "FAILED: %s (%v)\n", result.Path, result.Error)
		} else {
			fmt.Printf("OK: %s\n", result.Path)
		}
	}

	success, failed := eps.Summarize(results)
	fmt.Printf("Summary: success=%d failed=%d total=%d\n", success, failed, len(results))
	if failed > 0 {
		return errors.New("one or more files failed")
	}
	return nil
}

func parseOptions(args []string, allowFromDir bool) (options, []string, error) {
	opts := options{
		InputEncoding:  "gb2312",
		OutputEncoding: "utf-8",
	}
	positional := make([]string, 0)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--encoding":
			if i+1 >= len(args) {
				return opts, nil, errors.New("--encoding requires a value")
			}
			opts.InputEncoding = args[i+1]
			i++
		case "--to":
			if i+1 >= len(args) {
				return opts, nil, errors.New("--to requires a value")
			}
			opts.OutputEncoding = args[i+1]
			i++
		case "--from-dir":
			if !allowFromDir {
				return opts, nil, errors.New("--from-dir is only available in batch mode")
			}
			if i+1 >= len(args) {
				return opts, nil, errors.New("--from-dir requires a value")
			}
			opts.FromDir = args[i+1]
			i++
		case "--recursive":
			if !allowFromDir {
				return opts, nil, errors.New("--recursive is only available in batch mode")
			}
			opts.Recursive = true
		case "-h", "--help":
			opts.Help = true
		default:
			if strings.HasPrefix(arg, "-") {
				return opts, nil, fmt.Errorf("unknown flag: %s", arg)
			}
			positional = append(positional, filepath.Clean(arg))
		}
	}

	return opts, positional, nil
}

func printRootHelp() {
	fmt.Println(`eps - Minimal CSV cleaner

Usage:
  eps <input.csv> [output.csv] [--encoding <enc>] [--to <enc>]
  eps batch <file1.csv> <file2.csv> ... [--encoding <enc>] [--to <enc>]
  eps batch --from-dir <dir> [--recursive] [--encoding <enc>] [--to <enc>]

Use "eps --help", "eps batch --help" for more details.`)
}

func printSingleHelp() {
	fmt.Println(`Usage:
  eps <input.csv> [output.csv] [--encoding <enc>] [--to <enc>]

Defaults:
  --encoding gb2312
  --to       utf-8`)
}

func printBatchHelp() {
	fmt.Println(`Usage:
  eps batch <file1.csv> <file2.csv> ... [--encoding <enc>] [--to <enc>]
  eps batch --from-dir <dir> [--recursive] [--encoding <enc>] [--to <enc>]

Batch behavior:
  - Each file is processed independently.
  - Failures do not stop remaining files.
  - Outputs overwrite source files safely.`)
}
