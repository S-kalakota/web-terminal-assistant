package terminal

import (
	"os"
	"os/exec"
	"sync"
	"time"
)

// SessionOptions controls how a terminal session is started.
type SessionOptions struct {
	CWD   string
	Shell string
	Cols  int
	Rows  int
}

// Session owns one persistent shell process behind a PTY.
type Session struct {
	cmd         *exec.Cmd
	file        *os.File
	shell       string
	fallbackCWD string
	exited      chan struct{}
	waitErr     error
	waitMu      sync.Mutex
	closeOnce   sync.Once
	writeMu     sync.Mutex
}

// NewSession starts a persistent PTY-backed shell session.
func NewSession(opts SessionOptions) (*Session, error) {
	if opts.CWD == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		opts.CWD = cwd
	}
	if opts.Shell == "" {
		opts.Shell = DetectShell()
	}

	cmd, file, err := startPTY(opts)
	if err != nil {
		return nil, err
	}

	session := &Session{
		cmd:         cmd,
		file:        file,
		shell:       opts.Shell,
		fallbackCWD: opts.CWD,
		exited:      make(chan struct{}),
	}

	go func() {
		err := cmd.Wait()
		session.waitMu.Lock()
		session.waitErr = err
		session.waitMu.Unlock()
		close(session.exited)
	}()

	return session, nil
}

// Read reads raw PTY output.
func (s *Session) Read(buf []byte) (int, error) {
	return s.file.Read(buf)
}

// Write forwards browser input into the PTY.
func (s *Session) Write(data string) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	_, err := s.file.WriteString(data)
	return err
}

// Resize updates the PTY dimensions.
func (s *Session) Resize(cols int, rows int) error {
	if cols <= 0 || rows <= 0 {
		return nil
	}

	return resizePTY(s.file, cols, rows)
}

// Status reports the current best-effort shell state.
func (s *Session) Status() TerminalStatusMessage {
	running := true
	select {
	case <-s.exited:
		running = false
	default:
	}

	pid := 0
	if s.cmd.Process != nil {
		pid = s.cmd.Process.Pid
	}

	return TerminalStatusMessage{
		Type:    MessageTypeStatus,
		CWD:     currentWorkingDirectory(pid, s.fallbackCWD),
		Running: running,
		Shell:   s.shell,
	}
}

// Close tears down the PTY and shell process.
func (s *Session) Close() error {
	var err error
	s.closeOnce.Do(func() {
		if s.file != nil {
			_ = s.file.Close()
		}

		if s.cmd.Process == nil {
			return
		}

		select {
		case <-s.exited:
			err = s.waitError()
			return
		default:
		}

		_ = s.cmd.Process.Kill()

		select {
		case <-s.exited:
			err = s.waitError()
		case <-time.After(2 * time.Second):
		}
	})

	return err
}

func (s *Session) waitError() error {
	s.waitMu.Lock()
	defer s.waitMu.Unlock()

	return s.waitErr
}
