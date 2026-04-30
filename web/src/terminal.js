import { FitAddon } from "../node_modules/@xterm/addon-fit/lib/addon-fit.mjs";
import { Terminal } from "../node_modules/@xterm/xterm/lib/xterm.mjs";
import { encodeTerminalInput, encodeTerminalResize, openTerminalSocket } from "./api.js";
import { setConnectionStatus, setTerminalReady, setTerminalStatus } from "./state.js";

const TERMINAL_OPTIONS = {
  allowProposedApi: true,
  cursorBlink: true,
  cursorStyle: "bar",
  fontFamily: '"SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace',
  fontSize: 14,
  lineHeight: 1.24,
  scrollback: 8000,
  theme: {
    background: "#0b0d0e",
    foreground: "#e8edf0",
    cursor: "#9be7c2",
    cursorAccent: "#0b0d0e",
    selectionBackground: "#315b70",
    black: "#0b0d0e",
    red: "#ff7b72",
    green: "#9be7c2",
    yellow: "#f5c26b",
    blue: "#75b7ff",
    magenta: "#d7a6ff",
    cyan: "#7ddfd4",
    white: "#e8edf0",
    brightBlack: "#636b72",
    brightRed: "#ff9b94",
    brightGreen: "#b6f3d4",
    brightYellow: "#ffd98a",
    brightBlue: "#9fcbff",
    brightMagenta: "#e4c3ff",
    brightCyan: "#a5eee7",
    brightWhite: "#ffffff"
  }
};

export function createTerminalSession({ container, reconnectButton }) {
  const terminal = new Terminal(TERMINAL_OPTIONS);
  const fitAddon = new FitAddon();
  let socket = null;
  let resizeFrame = 0;

  terminal.loadAddon(fitAddon);
  terminal.open(container);
  setTerminalReady(true);

  terminal.writeln("Web Terminal");
  terminal.writeln("Connecting to local shell...\r\n");

  terminal.onData((data) => {
    if (!isSocketOpen(socket)) {
      return;
    }

    socket.send(encodeTerminalInput(data));
  });

  const resizeObserver = new ResizeObserver(() => {
    scheduleFitAndResize();
  });

  resizeObserver.observe(container);
  window.addEventListener("resize", scheduleFitAndResize);
  reconnectButton?.addEventListener("click", reconnect);

  queueMicrotask(() => {
    fitTerminal();
    connect();
  });

  function connect() {
    closeSocket();
    setConnectionStatus("connecting");
    reconnectButton.disabled = true;

    try {
      socket = openTerminalSocket();
    } catch (error) {
      handleConnectionError(error);
      return;
    }

    const activeSocket = socket;

    activeSocket.addEventListener("open", () => {
      if (socket !== activeSocket) {
        return;
      }

      terminal.writeln("\r\nConnected.\r\n");
      setConnectionStatus("connected");
      reconnectButton.disabled = false;
      sendResize();
    });

    activeSocket.addEventListener("message", (event) => {
      if (socket !== activeSocket) {
        return;
      }

      handleSocketMessage(event.data);
    });

    activeSocket.addEventListener("close", () => {
      if (socket !== activeSocket) {
        return;
      }

      setConnectionStatus("offline");
      reconnectButton.disabled = false;
    });

    activeSocket.addEventListener("error", () => {
      if (socket !== activeSocket) {
        return;
      }

      handleConnectionError(new Error("Terminal WebSocket connection failed."));
    });
  }

  function reconnect() {
    terminal.writeln("\r\nReconnecting...\r\n");
    connect();
  }

  function handleSocketMessage(data) {
    const message = parseMessage(data);
    if (!message) {
      terminal.write(String(data));
      return;
    }

    switch (message.type) {
      case "output":
        terminal.write(message.data ?? "");
        break;
      case "status":
        setTerminalStatus(message);
        break;
      case "error":
        terminal.writeln(`\r\n${message.message ?? "Terminal backend error."}\r\n`);
        setConnectionStatus("error", message.message ?? "");
        break;
      default:
        if (typeof message.data === "string") {
          terminal.write(message.data);
        }
    }
  }

  function handleConnectionError(error) {
    setConnectionStatus("error", error.message);
    reconnectButton.disabled = false;
    terminal.writeln(`\r\n${error.message}`);
    terminal.writeln("Start the Go server with terminal WebSocket support, then reconnect.\r\n");
  }

  function fitTerminal() {
    try {
      fitAddon.fit();
    } catch {
      return false;
    }

    return true;
  }

  function scheduleFitAndResize() {
    if (resizeFrame) {
      cancelAnimationFrame(resizeFrame);
    }

    resizeFrame = requestAnimationFrame(() => {
      resizeFrame = 0;
      if (fitTerminal()) {
        sendResize();
      }
    });
  }

  function sendResize() {
    if (!isSocketOpen(socket)) {
      return;
    }

    socket.send(encodeTerminalResize(terminal.cols, terminal.rows));
  }

  function sendCommand(command) {
    const normalized = command.endsWith("\n") || command.endsWith("\r") ? command : `${command}\r`;
    if (isSocketOpen(socket)) {
      socket.send(encodeTerminalInput(normalized));
    }
  }

  function closeSocket() {
    if (socket && socket.readyState < WebSocket.CLOSING) {
      socket.close();
    }
  }

  return {
    connect,
    dispose() {
      closeSocket();
      resizeObserver.disconnect();
      window.removeEventListener("resize", scheduleFitAndResize);
      terminal.dispose();
      setTerminalReady(false);
    },
    focus() {
      terminal.focus();
    },
    sendCommand
  };
}

function isSocketOpen(socket) {
  return socket?.readyState === WebSocket.OPEN;
}

function parseMessage(data) {
  if (typeof data !== "string") {
    return null;
  }

  try {
    return JSON.parse(data);
  } catch {
    return null;
  }
}
