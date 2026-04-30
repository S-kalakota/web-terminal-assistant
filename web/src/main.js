import { checkHealth, openTerminalSocket, suggestCommand } from "./api.js";
import { setConnectionStatus } from "./state.js";

const statusEl = document.querySelector("#connection-status");
const terminalOutputEl = document.querySelector("#terminal-output");
const connectButton = document.querySelector("#connect-terminal");
const assistantForm = document.querySelector("#assistant-form");
const assistantPrompt = document.querySelector("#assistant-prompt");
const assistantResult = document.querySelector("#assistant-result");

window.addEventListener("connection-change", (event) => {
  const connected = event.detail.connected;
  statusEl.textContent = connected ? "Server Ready" : "Disconnected";
  statusEl.dataset.state = connected ? "ready" : "offline";
});

async function boot() {
  try {
    await checkHealth();
    setConnectionStatus(true);
  } catch (error) {
    setConnectionStatus(false);
    appendTerminalLine(`health check failed: ${error.message}`);
  }
}

connectButton.addEventListener("click", () => {
  appendTerminalLine("opening placeholder websocket route...");
  const socket = openTerminalSocket();

  socket.addEventListener("open", () => {
    appendTerminalLine("websocket opened; terminal backend implementation is still pending");
  });

  socket.addEventListener("message", (event) => {
    appendTerminalLine(event.data);
  });

  socket.addEventListener("error", () => {
    appendTerminalLine("websocket connection failed because the skeleton route is not upgraded yet");
  });
});

assistantForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  assistantResult.textContent = "Thinking...";

  try {
    const result = await suggestCommand(assistantPrompt.value);
    renderSuggestions(result.suggestions ?? []);
  } catch (error) {
    assistantResult.textContent = error.message;
  }
});

function renderSuggestions(suggestions) {
  if (suggestions.length === 0) {
    assistantResult.textContent = "No suggestions yet.";
    return;
  }

  assistantResult.replaceChildren(
    ...suggestions.map((suggestion) => {
      const card = document.createElement("article");
      card.className = "suggestion";

      const command = document.createElement("code");
      command.textContent = suggestion.command;

      const explanation = document.createElement("p");
      explanation.textContent = suggestion.explanation;

      const risk = document.createElement("span");
      risk.className = `risk risk-${suggestion.risk}`;
      risk.textContent = `Risk: ${suggestion.risk}`;

      card.append(command, explanation, risk);
      return card;
    })
  );
}

function appendTerminalLine(line) {
  terminalOutputEl.textContent += `\n$ ${line}`;
}

boot();

