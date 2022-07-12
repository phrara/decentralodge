package tool

import (
	"os"
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
