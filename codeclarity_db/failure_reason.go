package codeclarity

import (
	"regexp"
	"strings"
)

// FailureReasonMaxLen matches the API's varchar(500) failure_reason column.
const FailureReasonMaxLen = 500

// StepFailureReasonKey is the key a plugin sets in its step result to say why it failed.
const StepFailureReasonKey = "failure_reason"

// Each pattern hides one kind of sensitive text. Order matters: credentials go before paths.
var failureReasonRedactions = []struct {
	pattern     *regexp.Regexp
	replacement string
}{
	{regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9+.-]*://)[^/\s@]+@`), "${1}"},
	{regexp.MustCompile(`(?i)\b(bearer\s+)[^\s"']+`), "${1}[secret]"},
	{regexp.MustCompile(`(?i)\b(password|passwd|secret|token|api[_-]?key)(\s*[=:]\s*)[^\s"',;]+`), "${1}${2}[secret]"},
	{regexp.MustCompile(`\b(ghp|gho|ghu|ghs|ghr|github_pat)_[A-Za-z0-9_]+`), "[secret]"},
	{regexp.MustCompile(`\bglpat-[A-Za-z0-9_-]+`), "[secret]"},
	{regexp.MustCompile(`(^|[\s"'(=])/[^\s"'),:;]+`), "${1}[path]"},
}

// SanitizeFailureReason makes a reason safe to display: it hides secrets and
// file paths, and cuts it to the column length without splitting a character.
func SanitizeFailureReason(reason string) string {
	for _, r := range failureReasonRedactions {
		reason = r.pattern.ReplaceAllString(reason, r.replacement)
	}
	reason = strings.TrimSpace(reason)
	if runes := []rune(reason); len(runes) > FailureReasonMaxLen {
		reason = string(runes[:FailureReasonMaxLen])
	}
	return reason
}
