package terminal

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/creack/pty"
)

const fallbackShell = "/bin/sh"

// DetectShell returns the user's configured shell with a conservative fallback.
func DetectShell() string {
	shell := strings.TrimSpace(os.Getenv("SHELL"))
	if shell == "" {
		return fallbackShell
	}

	return shell
}

func startPTY(opts SessionOptions) (*exec.Cmd, *os.File, error) {
	shell := opts.Shell
	if shell == "" {
		shell = DetectShell()
	}

	cmd := exec.Command(shell)
	if opts.CWD != "" {
		cmd.Dir = opts.CWD
	}
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	rows := uint16(opts.Rows)
	cols := uint16(opts.Cols)
	if rows == 0 {
		rows = 24
	}
	if cols == 0 {
		cols = 80
	}

	file, err := pty.StartWithSize(cmd, &pty.Winsize{
		Rows: rows,
		Cols: cols,
	})
	if err != nil {
		return nil, nil, err
	}

	return cmd, file, nil
}

func resizePTY(file *os.File, cols int, rows int) error {
	return pty.Setsize(file, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	})
}

func currentWorkingDirectory(pid int, fallback string) string {
	if pid <= 0 {
		return fallback
	}

	switch runtime.GOOS {
	case "linux":
		cwd, err := os.Readlink(filepath.Join("/proc", strconv.Itoa(pid), "cwd"))
		if err == nil && cwd != "" {
			return cwd
		}
	case "darwin":
		cwd := darwinProcessCWD(pid)
		if cwd != "" {
			return cwd
		}
	}

	return fallback
}

func darwinProcessCWD(pid int) string {
	output, err := exec.Command("lsof", "-a", "-d", "cwd", "-Fn", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(output), "\n") {
		if strings.HasPrefix(line, "n/") {
			return strings.TrimPrefix(line, "n")
		}
	}

	return ""
}
