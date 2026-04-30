package assistant

import (
	"regexp"
	"strings"
	"unicode"
)

var folderCreationPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^\s*(make|create)\s+(a\s+)?(folder|directory)\s+(named|called)\s+(.+?)\s*$`),
	regexp.MustCompile(`(?i)^\s*(make|create)\s+(a\s+)?(folder|directory)\s+(.+?)\s*$`),
}

func normalize(input string) string {
	return strings.Join(strings.Fields(strings.ToLower(input)), " ")
}

func matchKnownRequest(normalized string) (SuggestResponse, bool) {
	switch normalized {
	case "where am i", "print working directory", "what folder am i in", "current directory", "show current directory":
		return suggestion("pwd", "Prints the current working directory.", RiskLow), true
	case "list files", "show files", "show directory contents", "list directory contents":
		return suggestion("ls", "Lists files in the current directory.", RiskLow), true
	case "show hidden files", "list hidden files", "show all files":
		return suggestion("ls -la", "Lists files with details and includes hidden files.", RiskLow), true
	case "show files sorted by newest", "list files sorted by newest", "show newest files", "list newest files":
		return suggestion("ls -laht", "Lists all files with readable sizes and sorts newest first.", RiskLow), true
	case "what changed in git", "show git status", "git status", "what changed":
		return suggestion("git status", "Shows the current Git branch state and changed files.", RiskLow), true
	case "show current branch", "what branch am i on", "current git branch":
		return suggestion("git branch --show-current", "Prints the current Git branch name.", RiskLow), true
	case "find large files", "show large files":
		return suggestion("find . -type f -size +100M", "Finds files larger than 100 MB under the current directory.", RiskLow), true
	default:
		return SuggestResponse{}, false
	}
}

func matchFolderCreation(input string, normalized string) (SuggestResponse, bool) {
	if normalized == "make a folder" || normalized == "make folder" || normalized == "create a folder" || normalized == "create folder" ||
		normalized == "make a directory" || normalized == "make directory" || normalized == "create a directory" || normalized == "create directory" {
		return clarification("Which folder should I create?"), true
	}

	for _, pattern := range folderCreationPatterns {
		matches := pattern.FindStringSubmatch(input)
		if len(matches) == 0 {
			continue
		}

		name := strings.TrimSpace(matches[len(matches)-1])
		name = trimMatchingQuotes(name)
		if name == "" {
			return clarification("Which folder should I create?"), true
		}
		if looksUnsafeName(name) {
			return clarification("Please use a plain folder name without shell operators."), true
		}

		return suggestion("mkdir "+shellQuote(name), "Creates the requested folder in the current directory.", RiskLow), true
	}

	return SuggestResponse{}, false
}

func matchAmbiguousRequest(normalized string) (SuggestResponse, bool) {
	switch {
	case strings.Contains(normalized, "delete") || strings.Contains(normalized, "remove"):
		return clarification("What exactly should be removed? Destructive commands need a specific target and confirmation."), true
	case strings.Contains(normalized, "clean"):
		return clarification("What should be cleaned? I need a specific safe target before suggesting a command."), true
	case strings.Contains(normalized, "reset git") || normalized == "reset":
		return clarification("What kind of Git reset do you mean? Reset commands can discard work."), true
	case strings.Contains(normalized, "fix git"):
		return clarification("What Git problem should I inspect? A safe first step is asking for git status."), true
	case strings.Contains(normalized, "install"):
		return clarification("Which dependency or package manager should be used? Install commands can change this project."), true
	default:
		return SuggestResponse{}, false
	}
}

func trimMatchingQuotes(value string) string {
	if len(value) < 2 {
		return value
	}
	first := value[0]
	last := value[len(value)-1]
	if (first == '\'' && last == '\'') || (first == '"' && last == '"') {
		return strings.TrimSpace(value[1 : len(value)-1])
	}
	return value
}

func looksUnsafeName(value string) bool {
	return strings.ContainsAny(value, "\x00\n\r;&|`$<>")
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}
	if isSafeShellWord(value) {
		return value
	}
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}

func isSafeShellWord(value string) bool {
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			continue
		}
		switch r {
		case '.', '_', '-', '/', ':':
			continue
		default:
			return false
		}
	}
	return true
}
