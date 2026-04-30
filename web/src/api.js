const jsonHeaders = {
  "Content-Type": "application/json"
};

export async function checkHealth() {
  const response = await fetch("/healthz");
  if (!response.ok) {
    throw new Error(`Health check failed with ${response.status}`);
  }
  return response.json();
}

export async function suggestCommand(prompt, cwd = "") {
  const response = await fetch("/api/assistant/suggest", {
    method: "POST",
    headers: jsonHeaders,
    body: JSON.stringify({ prompt, cwd })
  });

  if (!response.ok) {
    throw new Error(`Suggestion request failed with ${response.status}`);
  }

  return response.json();
}

export async function checkCommandRisk(command) {
  const response = await fetch("/api/commands/risk", {
    method: "POST",
    headers: jsonHeaders,
    body: JSON.stringify({ command })
  });

  if (!response.ok) {
    throw new Error(`Risk request failed with ${response.status}`);
  }

  return response.json();
}

export async function recordAssistantCommand({ command, cwd = "", risk = "low" }) {
  const response = await fetch("/api/commands/audit", {
    method: "POST",
    headers: jsonHeaders,
    body: JSON.stringify({ command, cwd, risk, source: "assistant" })
  });

  if (!response.ok) {
    throw new Error(`Audit request failed with ${response.status}`);
  }

  return response.json();
}

export function openTerminalSocket() {
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  return new WebSocket(`${protocol}//${window.location.host}/ws/terminal`);
}

export function encodeTerminalInput(data) {
  return JSON.stringify({ type: "input", data });
}

export function encodeTerminalResize(cols, rows) {
  return JSON.stringify({ type: "resize", cols, rows });
}
