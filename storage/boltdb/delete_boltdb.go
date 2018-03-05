package boltdb

import (
	"fmt"
	"os"
	"path/filepath"
)

func DeleteBoltFile(path, name string) error {
	path = filepath.Join(path, name)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("deleting blockchain file: %s", err)
	}
	return nil
}
