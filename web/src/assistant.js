import { checkCommandRisk, recordAssistantCommand, suggestCommand } from "./api.js";
import { state } from "./state.js";

const RISK_LABELS = {
  low: "Low",
  medium: "Medium",
  high: "High"
};

export function createAssistantPanel({ container, runCommand }) {
  const view = {
    prompt: "",
    loading: false,
    message: "",
    messageTone: "neutral",
    suggestions: []
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
      view.message = view.suggestions.length ? "" : response.clarification || "No command suggestion available.";
      view.messageTone = "neutral";
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
    } catch (error) {
      view.message = error.message;
      view.messageTone = "danger";
    }

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

function escapeText(value) {
  return String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}

function escapeAttribute(value) {
  return escapeText(value).replaceAll('"', "&quot;");
}
