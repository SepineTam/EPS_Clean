# eps

`eps` is a minimal CSV cleaning CLI.

## What it does

For each input CSV file, `eps`:

1. Reads file bytes using an input encoding (default: `gb2312`)
2. Removes the last **three physical lines**
3. Writes the result using an output encoding (default: `utf-8`)

This tool works on plain text lines. It does not parse or reorder CSV fields.

## Installation

### Build from source

```bash
go build -o eps ./cmd/eps
```

### Install with Go

```bash
go install github.com/SepineTam/EPS_Clean/cmd/eps@latest
```

## Usage

```bash
eps <input.csv> [output.csv] [--encoding <input-encoding>] [--to <output-encoding>]
eps batch <file1.csv> <file2.csv> ... [--encoding <input-encoding>] [--to <output-encoding>]
eps batch --from-dir <dir> [--recursive] [--encoding <input-encoding>] [--to <output-encoding>]
```

## Examples

```bash
eps ori_file.csv
eps ori_file.csv new_file.csv
eps ori_file.csv new_file.csv --encoding gb2312 --to gbk
eps batch file1.csv file2.csv file3.csv --encoding gb2312 --to utf-8
eps batch --from-dir /path/to/a/set/of/files
```

## Supported encodings

Supported encoding names:

- `gb2312`
- `gbk`
- `utf-8`

Name matching is case-insensitive and tolerant of common forms such as `UTF8`, `utf_8`, `GBK`, and `cp936`.

### Encoding mapping note

`eps` uses the system `iconv` backend for non-UTF-8 conversion. In practice, `gb2312` and `gbk` are both converted through iconv-compatible code pages. If your platform maps these names differently, verify with a sample file on your system.

## Batch mode behavior

- `eps batch file1.csv file2.csv ...`: processes explicit files and overwrites each source safely.
- `eps batch --from-dir DIR`: scans `DIR` for `.csv` files (non-recursive by default).
- `--recursive`: also includes `.csv` files in subdirectories.
- Failures in one file do **not** stop other files.
- The CLI prints per-file status and a final summary.

## Notes

- Overwrite mode is safe: `eps` writes to a temp file first, then replaces the original file.
- Non-UTF-8 conversion depends on `iconv` being available in your runtime environment.
- If a file has **3 lines or fewer**, processing fails with a clear error.
- Common line endings are handled (`LF` and `CRLF`).

## Development

```bash
go test ./...
go build ./cmd/eps
```

## Release

GitHub Actions workflow builds release binaries for:

- macOS arm64
- macOS amd64
- Linux amd64
- Windows amd64

On tag push (for example `v1.0.0`), it uploads zipped binaries to GitHub Releases.
