import { checkCommandRisk, recordAssistantCommand, suggestCommand } from "./api.js";
import { state } from "./state.js";

const RISK_LABELS = {
  low: "Low",
  medium: "Medium",
  high: "High"
};

const HISTORY_LIMIT = 25;
const HISTORY_STORAGE_KEY = "web-terminal.assistant.history.v1";

export function createAssistantPanel({ container, runCommand }) {
  const view = {
    prompt: "",
    loading: false,
    message: "",
    messageTone: "neutral",
    suggestions: [],
    history: loadHistory(),
    activeHistoryId: ""
  };

  render();

  function render() {
    container.innerHTML = `
      <form class="assistant-form" data-role="assistant-form">
        <label class="assistant-label" for="assistant-prompt">Request</label>
        <textarea id="assistant-prompt" class="assistant-input" rows="4" placeholder="Show hidden files sorted by newest">${escapeText(view.prompt)}</textarea>
        <button class="assistant-submit" type="submit" ${view.loading ? "disabled" : ""}>${view.loading ? "Checking" : "Suggest"}</button>
      </form>
      <div class="assistant-message" data-tone="${view.messageTone}" ${view.message ? "" : "hidden"}>${escapeText(view.message)}</div>
      <div class="suggestion-list">
        ${view.suggestions.map(renderSuggestion).join("")}
      </div>
      ${renderHistory(view.history)}
    `;

    const form = container.querySelector('[data-role="assistant-form"]');
    const promptEl = container.querySelector("#assistant-prompt");

    form.addEventListener("submit", async (event) => {
      event.preventDefault();
      view.prompt = promptEl.value;
      await requestSuggestions();
    });

    for (const button of container.querySelectorAll("[data-run-index]")) {
      button.addEventListener("click", async () => {
        await runSuggestion(Number(button.dataset.runIndex));
      });
    }

    for (const input of container.querySelectorAll("[data-confirm-index]")) {
      input.addEventListener("input", () => {
        const suggestion = view.suggestions[Number(input.dataset.confirmIndex)];
        suggestion.confirmation = input.value;
      });
    }

    for (const button of container.querySelectorAll("[data-history-id]")) {
      button.addEventListener("click", () => {
        restoreHistoryEntry(button.dataset.historyId);
      });
    }

    const clearButton = container.querySelector("[data-clear-history]");
    clearButton?.addEventListener("click", () => {
      if (!window.confirm("Clear assistant history saved in this browser?")) {
        return;
      }
      view.history = [];
      view.activeHistoryId = "";
      saveHistory(view.history);
      render();
    });
  }

  async function requestSuggestions() {
    const prompt = view.prompt.trim();
    if (!prompt) {
      view.message = "Enter a request first.";
      view.messageTone = "warning";
      view.suggestions = [];
      render();
      return;
    }

    view.loading = true;
    view.message = "";
    view.suggestions = [];
    render();

    try {
      const response = await suggestCommand(prompt, state.cwd);
      const suggestions = response.suggestions ?? [];
      view.suggestions = await Promise.all(suggestions.map(enrichSuggestion));
      view.message = response.warning || (view.suggestions.length ? "" : response.clarification || "No command suggestion available.");
      view.messageTone = response.warning ? "warning" : "neutral";
      view.activeHistoryId = addHistoryEntry({
        prompt,
        cwd: state.cwd,
        suggestions: view.suggestions,
        clarification: response.clarification || "",
        status: view.suggestions.length ? "suggested" : "clarification"
      });
    } catch (error) {
      view.message = error.message;
      view.messageTone = "danger";
    } finally {
      view.loading = false;
      render();
    }
  }

  async function enrichSuggestion(suggestion) {
    const assessment = await checkCommandRisk(suggestion.command);
    return {
      ...suggestion,
      risk: assessment.risk ?? suggestion.risk ?? "low",
      reason: assessment.reason ?? "",
      requiresConfirmation: Boolean(assessment.requiresConfirmation),
      confirmation: ""
    };
  }

  async function runSuggestion(index) {
    const suggestion = view.suggestions[index];
    if (!suggestion) {
      return;
    }

    if (suggestion.risk === "high" && suggestion.confirmation !== suggestion.command) {
      view.message = "Type the exact command before running a high-risk suggestion.";
      view.messageTone = "danger";
      render();
      return;
    }

    try {
      await recordAssistantCommand({
        command: suggestion.command,
        cwd: state.cwd,
        risk: suggestion.risk
      });
      runCommand(suggestion.command);
      view.message = "Command sent to terminal.";
      view.messageTone = "success";
      markHistoryEntryRun(suggestion);
    } catch (error) {
      view.message = error.message;
      view.messageTone = "danger";
    }

    render();
  }

  function addHistoryEntry(entry) {
    const id = `${Date.now()}-${Math.random().toString(16).slice(2)}`;
    const savedEntry = {
      id,
      prompt: entry.prompt,
      cwd: entry.cwd || "",
      suggestions: sanitizeSuggestions(entry.suggestions),
      clarification: entry.clarification || "",
      status: entry.status,
      createdAt: new Date().toISOString(),
      ranCommand: "",
      ranAt: ""
    };

    view.history = [savedEntry, ...view.history.filter((item) => item.prompt !== savedEntry.prompt)].slice(
      0,
      HISTORY_LIMIT
    );
    saveHistory(view.history);
    return id;
  }

  function markHistoryEntryRun(suggestion) {
    if (!view.activeHistoryId) {
      return;
    }

    view.history = view.history.map((entry) => {
      if (entry.id !== view.activeHistoryId) {
        return entry;
      }

      return {
        ...entry,
        status: "ran",
        ranCommand: suggestion.command,
        ranAt: new Date().toISOString()
      };
    });
    saveHistory(view.history);
  }

  function restoreHistoryEntry(id) {
    const entry = view.history.find((item) => item.id === id);
    if (!entry) {
      return;
    }

    view.prompt = entry.prompt;
    view.suggestions = (entry.suggestions || []).map((suggestion) => ({
      ...suggestion,
      confirmation: ""
    }));
    view.activeHistoryId = entry.id;
    view.message = entry.clarification || "";
    view.messageTone = entry.clarification ? "neutral" : "success";
    render();
  }
}

function renderSuggestion(suggestion, index) {
  const risk = suggestion.risk ?? "low";
  const isHighRisk = risk === "high";
  const reason = suggestion.reason || suggestion.explanation || "";

  return `
    <article class="suggestion-card" data-risk="${escapeAttribute(risk)}">
      <div class="suggestion-card__header">
        <code>${escapeText(suggestion.command)}</code>
        <span class="risk-badge" data-risk="${escapeAttribute(risk)}">${escapeText(RISK_LABELS[risk] ?? risk)}</span>
      </div>
      <p>${escapeText(suggestion.explanation ?? "")}</p>
      ${risk === "medium" || risk === "high" ? `<div class="risk-warning">${escapeText(reason)}</div>` : ""}
      ${
        isHighRisk
          ? `<label class="strong-confirm">Confirm<code>${escapeText(suggestion.command)}</code><input type="text" autocomplete="off" spellcheck="false" data-confirm-index="${index}" value="${escapeAttribute(suggestion.confirmation ?? "")}" /></label>`
          : ""
      }
      <button type="button" data-run-index="${index}">Run</button>
    </article>
  `;
}

function renderHistory(history) {
  const items = history
    .map((entry) => {
      const command = entry.ranCommand || entry.suggestions?.[0]?.command || entry.clarification || "No command";
      const status = entry.status === "ran" ? "Ran" : entry.status === "clarification" ? "Clarified" : "Suggested";
      return `
        <li class="assistant-history-item">
          <button type="button" data-history-id="${escapeAttribute(entry.id)}">
            <span>${escapeText(entry.prompt)}</span>
            <small>${escapeText(status)} · ${escapeText(command)}</small>
          </button>
        </li>
      `;
    })
    .join("");

  return `
    <section class="assistant-history" aria-label="Assistant history">
      <div class="assistant-history__header">
        <h3>History</h3>
        <button type="button" data-clear-history ${history.length ? "" : "disabled"}>Clear</button>
      </div>
      ${
        history.length
          ? `<ol>${items}</ol>`
          : `<p class="assistant-history-empty">Assistant requests will be saved here on this browser.</p>`
      }
    </section>
  `;
}

function loadHistory() {
  try {
    const parsed = JSON.parse(window.localStorage.getItem(HISTORY_STORAGE_KEY) || "[]");
    if (!Array.isArray(parsed)) {
      return [];
    }
    return parsed.slice(0, HISTORY_LIMIT).map(normalizeHistoryEntry).filter(Boolean);
  } catch {
    return [];
  }
}

function saveHistory(history) {
  try {
    window.localStorage.setItem(HISTORY_STORAGE_KEY, JSON.stringify(history.slice(0, HISTORY_LIMIT)));
  } catch {
    // The assistant should keep working even if browser storage is unavailable.
  }
}

function normalizeHistoryEntry(entry) {
  if (!entry || typeof entry.prompt !== "string" || !entry.prompt.trim()) {
    return null;
  }

  return {
    id: typeof entry.id === "string" && entry.id ? entry.id : `${Date.now()}-${Math.random()}`,
    prompt: entry.prompt,
    cwd: typeof entry.cwd === "string" ? entry.cwd : "",
    suggestions: sanitizeSuggestions(entry.suggestions),
    clarification: typeof entry.clarification === "string" ? entry.clarification : "",
    status: typeof entry.status === "string" ? entry.status : "suggested",
    createdAt: typeof entry.createdAt === "string" ? entry.createdAt : new Date().toISOString(),
    ranCommand: typeof entry.ranCommand === "string" ? entry.ranCommand : "",
    ranAt: typeof entry.ranAt === "string" ? entry.ranAt : ""
  };
}

function sanitizeSuggestions(suggestions) {
  if (!Array.isArray(suggestions)) {
    return [];
  }

  return suggestions.slice(0, 3).map((suggestion) => ({
    command: String(suggestion.command || ""),
    explanation: String(suggestion.explanation || ""),
    risk: String(suggestion.risk || "low"),
    reason: String(suggestion.reason || ""),
    requiresConfirmation: Boolean(suggestion.requiresConfirmation),
    confirmation: ""
  }));
}

function escapeText(value) {
  return String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}

function escapeAttribute(value) {
  return escapeText(value).replaceAll('"', "&quot;");
}
