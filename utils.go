package heidou

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/tools/imports"
)

// format
func format(filename string, content []byte) error {
	ext := filepath.Ext(filename)
	data := content
	if ext == ".go" {
		var err error
		data, err = imports.Process(filename, content, nil)
		if err != nil {
			return fmt.Errorf("format file %s: %v", filename, err)
		}

	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("write file %s: %v", filename, err)
	}
	return nil
}
