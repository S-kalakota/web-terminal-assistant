package server

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"web-terminal/internal/terminal"
)

var terminalUpgrader = websocket.Upgrader{
	CheckOrigin: allowLocalWebSocketOrigin,
}

func (s *Server) handleTerminalWebSocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	conn, err := terminalUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	session, err := terminal.NewSession(terminal.SessionOptions{
		Cols: 80,
		Rows: 24,
	})
	if err != nil {
		_ = conn.WriteJSON(terminal.ErrorMessage{
			Type:    terminal.MessageTypeError,
			Message: "failed to start terminal session",
		})
		return
	}
	defer session.Close()

	var writeMu sync.Mutex
	writeJSON := func(payload any) error {
		writeMu.Lock()
		defer writeMu.Unlock()
		return conn.WriteJSON(payload)
	}

	if err := writeJSON(session.Status()); err != nil {
		return
	}

	done := make(chan struct{})
	defer close(done)

	errs := make(chan error, 3)

	go streamPTYOutput(session, writeJSON, errs)
	go receiveTerminalInput(conn, session, writeJSON, errs)
	go streamSessionStatus(session, writeJSON, done, errs)

	<-errs
}

func streamPTYOutput(session *terminal.Session, writeJSON func(any) error, errs chan<- error) {
	buf := make([]byte, 4096)
	for {
		n, err := session.Read(buf)
		if n > 0 {
			output := terminal.TerminalOutputMessage{
				Type: terminal.MessageTypeOutput,
				Data: string(buf[:n]),
			}
			if writeErr := writeJSON(output); writeErr != nil {
				errs <- writeErr
				return
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				errs <- nil
				return
			}
			errs <- err
			return
		}
	}
}

func receiveTerminalInput(conn *websocket.Conn, session *terminal.Session, writeJSON func(any) error, errs chan<- error) {
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			errs <- err
			return
		}

		msg, err := terminal.ParseClientMessage(payload)
		if err != nil {
			_ = writeJSON(terminal.ErrorMessage{
				Type:    terminal.MessageTypeError,
				Message: err.Error(),
			})
			continue
		}

		switch msg.Type {
		case terminal.MessageTypeInput:
			if err := session.Write(msg.Data); err != nil {
				errs <- err
				return
			}
		case terminal.MessageTypeResize:
			if err := session.Resize(msg.Cols, msg.Rows); err != nil {
				errs <- err
				return
			}
			if err := writeJSON(session.Status()); err != nil {
				errs <- err
				return
			}
		}
	}
}

func streamSessionStatus(session *terminal.Session, writeJSON func(any) error, done <-chan struct{}, errs chan<- error) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if err := writeJSON(session.Status()); err != nil {
				errs <- err
				return
			}
		}
	}
}

func allowLocalWebSocketOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}

	originURL, err := url.Parse(origin)
	if err != nil {
		return false
	}

	if strings.EqualFold(originURL.Host, r.Host) {
		return true
	}

	requestHost, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		requestHost = r.Host
	}

	return isLoopbackHost(originURL.Hostname()) && isLoopbackHost(requestHost)
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}

	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
