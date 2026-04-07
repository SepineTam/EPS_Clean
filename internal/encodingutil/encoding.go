package encodingutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"unicode/utf8"
)

func normalize(name string) string {
	s := strings.TrimSpace(strings.ToLower(name))
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "-", "")
	return s
}

func resolveName(name string) (string, bool, error) {
	switch normalize(name) {
	case "utf8", "utf":
		return "UTF-8", true, nil
	case "gb2312":
		return "GB2312", false, nil
	case "gbk", "cp936", "ms936":
		return "GBK", false, nil
	default:
		return "", false, fmt.Errorf("unsupported encoding: %q (supported: gb2312, gbk, utf-8)", name)
	}
}

func Decode(input []byte, encName string) (string, error) {
	source, utf8Mode, err := resolveName(encName)
	if err != nil {
		return "", err
	}
	if utf8Mode {
		if !utf8.Valid(input) {
			return "", fmt.Errorf("input is not valid UTF-8")
		}
		return string(input), nil
	}
	decoded, err := runIconv(input, source, "UTF-8")
	if err != nil {
		return "", fmt.Errorf("decode failed with %s: %w", encName, err)
	}
	return string(decoded), nil
}

func Encode(text string, encName string) ([]byte, error) {
	target, utf8Mode, err := resolveName(encName)
	if err != nil {
		return nil, err
	}
	if utf8Mode {
		return []byte(text), nil
	}
	encoded, err := runIconv([]byte(text), "UTF-8", target)
	if err != nil {
		return nil, fmt.Errorf("encode failed with %s: %w", encName, err)
	}
	return encoded, nil
}

func runIconv(input []byte, from, to string) ([]byte, error) {
	cmd := exec.Command("iconv", "-f", from, "-t", to)
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("iconv conversion %s->%s failed: %s", from, to, msg)
	}
	return stdout.Bytes(), nil
}
