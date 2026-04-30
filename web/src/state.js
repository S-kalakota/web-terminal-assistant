export const state = {
  cwd: "",
  connected: false,
  shell: "",
  terminalReady: false
};

export function setConnectionStatus(connected) {
  state.connected = connected;
  window.dispatchEvent(new CustomEvent("connection-change", { detail: { connected } }));
}

