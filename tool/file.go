package tool

import (
	"os"
	"strings"
)

func WriteFile(b []byte, path string) error {
	err := os.WriteFile(path, b, 02)
	if err != nil {
		return err
	}
	return nil
}

func LoadFile(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func WrapFile(content string) string {
	return "<content:" + content + ">"
}

func UnwrapFile(file string) string {
	content := strings.ReplaceAll(file, "<content:", "")
	content = strings.ReplaceAll(content, ">", "")
	return content
}
