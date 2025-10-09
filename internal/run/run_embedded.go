package run

import (
	"os"
	"path/filepath"

	"github.com/Cyrof/machina/internal/resources"
)

func PS1Embedded(name string, argv ...string) error {
	content, err := resources.ReadPS1(name)
	if err != nil {
		return err
	}

	tmpDir := os.TempDir()
	tmpPath := filepath.Join(tmpDir, "machina-"+name)
	if err := os.WriteFile(tmpPath, []byte(content), 0600); err != nil {
		return err
	}
	defer os.Remove(tmpPath)

	return PS1(tmpPath, argv...)
}
