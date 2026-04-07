package eps

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SepineTam/EPS_Clean/internal/encodingutil"
)

func writeEncodedFile(t *testing.T, path, content, enc string) {
	t.Helper()
	b, err := encodingutil.Encode(content, enc)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
}

func readDecodedFile(t *testing.T, path, enc string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	s, err := encodingutil.Decode(b, enc)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	return s
}

func TestProcessFile_GB2312ToUTF8(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.csv")
	out := filepath.Join(dir, "out.csv")
	writeEncodedFile(t, in, "列1\n1\n2\n3\n4\n", "gb2312")

	err := ProcessFile(FileTask{
		InputPath:      in,
		OutputPath:     out,
		InputEncoding:  "gb2312",
		OutputEncoding: "utf-8",
	})
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}

	got := readDecodedFile(t, out, "utf-8")
	if got != "列1\n1\n" {
		t.Fatalf("unexpected content: %q", got)
	}
}

func TestProcessFile_UTF8ToGBK(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.csv")
	out := filepath.Join(dir, "out.csv")
	writeEncodedFile(t, in, "h\n1\n2\n3\n4\n", "utf-8")

	err := ProcessFile(FileTask{
		InputPath:      in,
		OutputPath:     out,
		InputEncoding:  "utf-8",
		OutputEncoding: "gbk",
	})
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}

	got := readDecodedFile(t, out, "gbk")
	if got != "h\n1\n" {
		t.Fatalf("unexpected content: %q", got)
	}
}

func TestRemoveLastThreeLines_CRLF(t *testing.T) {
	result, err := removeLastNPhysicalLines("a\r\nb\r\nc\r\nd\r\ne\r\n", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "a\r\nb\r\n" {
		t.Fatalf("unexpected result: %q", result)
	}
}

func TestProcessFile_TooFewLines(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.csv")
	out := filepath.Join(dir, "out.csv")
	writeEncodedFile(t, in, "a\n1\n2\n", "utf-8")

	err := ProcessFile(FileTask{
		InputPath:      in,
		OutputPath:     out,
		InputEncoding:  "utf-8",
		OutputEncoding: "utf-8",
	})
	if err == nil {
		t.Fatal("expected error for too few lines")
	}
}

func TestProcessFile_SafeOverwrite(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "same.csv")
	writeEncodedFile(t, in, "a\n1\n2\n3\n4\n", "utf-8")

	err := ProcessFile(FileTask{
		InputPath:      in,
		OutputPath:     in,
		InputEncoding:  "utf-8",
		OutputEncoding: "utf-8",
		SafeOverwrite:  true,
	})
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}

	got := readDecodedFile(t, in, "utf-8")
	if got != "a\n1\n" {
		t.Fatalf("unexpected content: %q", got)
	}
}

func TestProcessFile_OutputToNewFile(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.csv")
	out := filepath.Join(dir, "new.csv")
	orig := "a\n1\n2\n3\n4\n"
	writeEncodedFile(t, in, orig, "utf-8")

	err := ProcessFile(FileTask{
		InputPath:      in,
		OutputPath:     out,
		InputEncoding:  "utf-8",
		OutputEncoding: "utf-8",
	})
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}

	if got := readDecodedFile(t, in, "utf-8"); got != orig {
		t.Fatalf("input file was modified: %q", got)
	}
	if got := readDecodedFile(t, out, "utf-8"); got != "a\n1\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestProcessBatch_PartialSuccess(t *testing.T) {
	dir := t.TempDir()
	okPath := filepath.Join(dir, "ok.csv")
	writeEncodedFile(t, okPath, "a\n1\n2\n3\n4\n", "utf-8")
	failPath := filepath.Join(dir, "missing.csv")

	results := ProcessBatch([]string{okPath, failPath}, BatchConfig{InputEncoding: "utf-8", OutputEncoding: "utf-8"})
	success, failed := Summarize(results)
	if success != 1 || failed != 1 {
		t.Fatalf("unexpected summary success=%d failed=%d", success, failed)
	}
}

func TestScanCSVFiles_FromDir(t *testing.T) {
	dir := t.TempDir()
	csvA := filepath.Join(dir, "a.csv")
	csvB := filepath.Join(dir, "b.CSV")
	txt := filepath.Join(dir, "c.txt")
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	nested := filepath.Join(sub, "d.csv")

	for _, p := range []string{csvA, csvB, txt, nested} {
		if err := os.WriteFile(p, []byte("x\n1\n2\n3\n4\n"), 0o644); err != nil {
			t.Fatalf("write failed: %v", err)
		}
	}

	nonRecursive, err := ScanCSVFiles(dir, false)
	if err != nil {
		t.Fatalf("ScanCSVFiles failed: %v", err)
	}
	if len(nonRecursive) != 2 {
		t.Fatalf("expected 2 csv files, got %d", len(nonRecursive))
	}

	recursive, err := ScanCSVFiles(dir, true)
	if err != nil {
		t.Fatalf("ScanCSVFiles recursive failed: %v", err)
	}
	if len(recursive) != 3 {
		t.Fatalf("expected 3 csv files, got %d", len(recursive))
	}
}
