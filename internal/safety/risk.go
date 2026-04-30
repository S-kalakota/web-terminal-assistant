package safety

import (
	"regexp"
	"strings"

	"web-terminal/internal/terminal"
)

// Assessment is the server-side risk decision for a command.
type Assessment struct {
	Risk   terminal.RiskLevel
	Reason string
}

var (
	remotePipePattern = regexp.MustCompile(`(?i)\b(curl|wget)\b.+\|\s*(sh|bash|zsh|python|ruby|perl|node)\b`)
	rmRecursive       = regexp.MustCompile(`(?i)\brm\s+[^;&|]*-[^\s;&|]*r[^\s;&|]*f|(?i)\brm\s+[^;&|]*-[^\s;&|]*f[^\s;&|]*r`)
)

// Classify assigns a conservative risk level to a shell command.
func Classify(command string) terminal.CommandRiskResponse {
	assessment := AssessCommand(command)
	return terminal.CommandRiskResponse{
		Risk:                 assessment.Risk,
		RequiresConfirmation: RequiresConfirmation(assessment.Risk),
		Reason:               assessment.Reason,
	}
}

// AssessCommand assigns a conservative risk level to a shell command.
func AssessCommand(command string) Assessment {
	normalized := strings.ToLower(strings.TrimSpace(command))
	if normalized == "" {
		return Assessment{
			Risk:   terminal.RiskLow,
			Reason: "Empty commands do not run anything.",
		}
	}

	switch {
	case isHighRisk(normalized):
		return Assessment{
			Risk:   terminal.RiskHigh,
			Reason: "This command can make destructive, privileged, or system-level changes.",
		}
	case isMediumRisk(normalized):
		return Assessment{
			Risk:   terminal.RiskMedium,
			Reason: "This command may change files, install software, or contact the network.",
		}
	default:
		return Assessment{
			Risk:   terminal.RiskLow,
			Reason: "This command appears read-only or low impact.",
		}
	}
}

func isHighRisk(command string) bool {
	highTokens := []string{
		"sudo ",
		"chmod -r",
		"chmod --recursive",
		"chown -r",
		"chown --recursive",
		"mkfs",
		"diskutil erasedisk",
		"format ",
		"dd if=",
		" /etc/",
		" /usr/bin/",
		" /bin/",
		" /sbin/",
	}

	if rmRecursive.MatchString(command) || remotePipePattern.MatchString(command) {
		return true
	}

	for _, token := range highTokens {
		if strings.Contains(command, token) {
			return true
		}
	}

	return false
}

func isMediumRisk(command string) bool {
	mediumTokens := []string{
		">",
		" mv ",
		"mv ",
		" cp ",
		"cp ",
		"npm install",
		"yarn add",
		"pnpm add",
		"pip install",
		"brew install",
		"apt install",
		"apt-get install",
		"curl ",
		"wget ",
		"mkdir ",
		"touch ",
	}

	for _, token := range mediumTokens {
		if strings.Contains(command, token) {
			return true
		}
	}

	return false
}
