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

var changeDirectoryPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^\s*(go|change|cd|move)\s+(to|into)\s+(the\s+)?(folder|directory)\s+(.+?)\s*$`),
	regexp.MustCompile(`(?i)^\s*(go|change|cd|move)\s+(to|into)\s+(.+?)\s*$`),
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
	case "go back a directory", "go back one directory", "go up a directory", "go up one directory", "back a directory", "up one directory", "parent directory", "change to parent directory", "move back a directory":
		return suggestion("cd ..", "Changes the terminal session to the parent directory.", RiskLow), true
	case "go home", "go to home", "change to home directory", "home directory":
		return suggestion("cd ~", "Changes the terminal session to your home directory.", RiskLow), true
	case "clear terminal", "clear the terminal", "clear screen", "clear the screen":
		return suggestion("clear", "Clears the terminal display.", RiskLow), true
	case "what time is it", "show date", "show time", "date and time":
		return suggestion("date", "Prints the current date and time.", RiskLow), true
	case "who am i", "show current user", "current user":
		return suggestion("whoami", "Prints the current username.", RiskLow), true
	case "show disk space", "disk space":
		return suggestion("df -h", "Shows mounted disk usage in human-readable units.", RiskLow), true
	case "show folder size", "current folder size", "directory size":
		return suggestion("du -sh .", "Shows the approximate size of the current directory.", RiskLow), true
	case "run tests", "test project":
		return suggestion("go test ./...", "Runs the Go test suite for this project.", RiskLow), true
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

func matchDirectoryNavigation(input string, normalized string) (SuggestResponse, bool) {
	for _, pattern := range changeDirectoryPatterns {
		matches := pattern.FindStringSubmatch(input)
		if len(matches) == 0 {
			continue
		}

		target := strings.TrimSpace(matches[len(matches)-1])
		target = trimMatchingQuotes(target)
		if target == "" {
			return clarification("Which directory should I change into?"), true
		}
		if looksUnsafeName(target) {
			return clarification("Please use a plain directory name without shell operators."), true
		}
		if target == "parent" || strings.EqualFold(target, "parent directory") {
			return suggestion("cd ..", "Changes the terminal session to the parent directory.", RiskLow), true
		}

		return suggestion("cd "+shellQuote(target), "Changes the terminal session to the requested directory.", RiskLow), true
	}

	if strings.Contains(normalized, "go back") || strings.Contains(normalized, "go up") {
		return suggestion("cd ..", "Changes the terminal session to the parent directory.", RiskLow), true
	}

	return SuggestResponse{}, false
}

func matchAmbiguousRequest(normalized string) (SuggestResponse, bool) {
	switch {
	case normalized == "move directory" || normalized == "move folder" || normalized == "move a directory" || normalized == "move a folder":
		return suggestion("mv <source_directory> <destination_directory>", "Moves or renames a directory. Replace both placeholders with the real paths.", RiskMedium), true
	case strings.Contains(normalized, "move") && (strings.Contains(normalized, "directory") || strings.Contains(normalized, "folder")):
		return clarification("Which directory should be moved, and where should it go?"), true
	case strings.Contains(normalized, "rename") && (strings.Contains(normalized, "directory") || strings.Contains(normalized, "folder")):
		return clarification("Which directory should be renamed, and what should the new name be?"), true
	case strings.Contains(normalized, "copy") && (strings.Contains(normalized, "directory") || strings.Contains(normalized, "folder")):
		return clarification("Which directory should be copied, and where should the copy go?"), true
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

func fallbackSuggestion(normalized string) (SuggestResponse, bool) {
	switch {
	case strings.Contains(normalized, "list") || strings.Contains(normalized, "show"):
		return suggestion("ls", "Lists files in the current directory.", RiskLow), true
	case strings.Contains(normalized, "git"):
		return suggestion("git status", "Shows the current Git branch state and changed files.", RiskLow), true
	case strings.Contains(normalized, "find") || strings.Contains(normalized, "search"):
		return suggestion("find . -name '<pattern>'", "Searches for files by name. Replace the placeholder with the pattern to match.", RiskLow), true
	case strings.Contains(normalized, "create") || strings.Contains(normalized, "make"):
		return suggestion("mkdir <name>", "Creates a directory. Replace the placeholder with the directory name.", RiskMedium), true
	case strings.Contains(normalized, "move") || strings.Contains(normalized, "rename"):
		return suggestion("mv <source> <destination>", "Moves or renames a file or directory. Replace both placeholders with real paths.", RiskMedium), true
	case strings.Contains(normalized, "copy"):
		return suggestion("cp -R <source> <destination>", "Copies a file or directory. Replace both placeholders with real paths.", RiskMedium), true
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
