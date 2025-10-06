package elevate

import (
	"os"
	"os/exec"
	"strings"
)

func IsAdmin() bool {
	cmd := exec.Command("cmd.exe", "/c", "net", "session")
	return cmd.Run() == nil
}

func RelaunchElevated() error {
	exe, _ := os.Executable()

	args := make([]string, 0, len(os.Args))
	for _, a := range os.Args[1:] {
		if a == "-elevated" || a == "--elevated" {
			continue
		}
		args = append(args, a)
	}
	args = append(args, "-elevated")

	ps := `Start-Process -FilePath "` + exe + `" -ArgumentList ` + psArgList(args) + ` -Verb RunAs`
	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", ps)
	return cmd.Run()
}

func psArgList(args []string) string {
	quoted := make([]string, 0, len(args))
	for _, a := range args {
		a = strings.ReplaceAll(a, "`", "``")
		quoted = append(quoted, "`\""+a+"`\"")
	}
	return strings.Join(quoted, ",")
}
