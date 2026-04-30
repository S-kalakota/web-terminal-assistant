import { checkHealth } from "./api.js";
import { createAssistantPanel } from "./assistant.js";
import { createTerminalSession } from "./terminal.js";
import { setConnectionStatus, state } from "./state.js";

const statusEl = document.querySelector("#connection-status");
const cwdEl = document.querySelector("#cwd");
const terminalEl = document.querySelector("#terminal");
const assistantEl = document.querySelector("#assistant-panel");
const reconnectButton = document.querySelector("#reconnect-terminal");
let terminalSession = null;

window.addEventListener("connection-change", (event) => {
  renderConnectionStatus(event.detail);
});

window.addEventListener("terminal-status", (event) => {
  renderConnectionStatus(event.detail);
  renderCwd(event.detail.cwd, event.detail.shell);
});

async function boot() {
  renderConnectionStatus(state);
  renderCwd(state.cwd, state.shell);

  try {
    await checkHealth();
    setConnectionStatus("connecting");
  } catch (error) {
    setConnectionStatus("offline", error.message);
  }

  terminalSession = createTerminalSession({
    container: terminalEl,
    reconnectButton
  });

  createAssistantPanel({
    container: assistantEl,
    runCommand(command) {
      terminalSession.sendCommand(command);
      terminalSession.focus();
    }
  });

  window.webTerminal = {
    runCommand(command) {
      terminalSession.sendCommand(command);
      terminalSession.focus();
    }
  };
}

function renderConnectionStatus({ connectionStatus, message }) {
  const labels = {
    checking: "Checking",
    connecting: "Connecting",
    connected: "Connected",
    error: "Connection Error",
    offline: "Disconnected"
  };

  statusEl.textContent = labels[connectionStatus] ?? "Disconnected";
  statusEl.dataset.state = connectionStatus ?? "offline";
  statusEl.title = message || statusEl.textContent;
}

function renderCwd(cwd, shell) {
  const parts = [cwd || "Waiting for shell", shell].filter(Boolean);
  cwdEl.textContent = parts.join("  ");
}

boot();
