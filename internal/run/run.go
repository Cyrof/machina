package run

import (
	"os"
	"os/exec"
)

func BAT(path string, argv ...string) error {
	args := append([]string{"/c", path}, argv...)
	cmd := exec.Command("cmd.exe", args...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd.Run()
}

func PS1(path string, argv ...string) error {
	args := []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-File", path}
	args = append(args, argv...)
	cmd := exec.Command("powershell.exe", args...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd.Run()
}