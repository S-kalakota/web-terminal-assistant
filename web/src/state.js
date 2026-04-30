export const state = {
  cwd: "",
  connectionStatus: "checking",
  connected: false,
  message: "",
  shell: "",
  terminalReady: false
};

export function setConnectionStatus(status, message = "") {
  if (typeof status === "boolean") {
    state.connectionStatus = status ? "connected" : "offline";
    state.connected = status;
  } else {
    state.connectionStatus = status;
    state.connected = status === "connected";
  }

  state.message = message;
  window.dispatchEvent(new CustomEvent("connection-change", { detail: { ...state } }));
}

export function setTerminalReady(terminalReady) {
  state.terminalReady = terminalReady;
  window.dispatchEvent(new CustomEvent("terminal-ready", { detail: { ...state } }));
}

export function setTerminalStatus({ cwd = state.cwd, running = state.connected, shell = state.shell }) {
  state.cwd = cwd || state.cwd;
  state.connected = Boolean(running);
  state.connectionStatus = running ? "connected" : "offline";
  state.shell = shell || state.shell;

  window.dispatchEvent(new CustomEvent("terminal-status", { detail: { ...state } }));
}
