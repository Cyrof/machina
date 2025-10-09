package resources

import (
	"embed"
	"io/fs"
)

// go:embed ps1/*.ps1
var ps1FS embed.FS

func ReadPS1(name string) (string, error) {
	b, err := fs.ReadFile(ps1FS, "ps1/"+name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
