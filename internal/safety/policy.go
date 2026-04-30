package safety

import "web-terminal/internal/terminal"

// RequiresConfirmation reports whether a command needs explicit user confirmation.
func RequiresConfirmation(risk terminal.RiskLevel) bool {
	return risk == terminal.RiskMedium || risk == terminal.RiskHigh
}

// RequiresStrongConfirmation reports whether the UI should ask for an explicit typed confirmation.
func RequiresStrongConfirmation(risk terminal.RiskLevel) bool {
	return risk == terminal.RiskHigh
}
