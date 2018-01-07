package boltdb

import (
	"os"
	"fmt"
)

func DeleteBlockchain(name string) error {
	path := fmt.Sprintf(dbFile, name)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("deleting blockchain file: %s", err)
	}
	return nil
}
