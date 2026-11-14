package codeclarity

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestSanitizeFailureReason(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{"keeps a plain reason", "no supported manifest was found", "no supported manifest was found"},
		{"hides a file path", "cannot read /private/org/projects/x/package.json now", "cannot read [path] now"},
		{"hides a quoted path", `open "/etc/passwd": denied`, `open "[path]": denied`},
		{"leaves words with slashes", "use npm/yarn and/or composer", "use npm/yarn and/or composer"},
		{"removes URL credentials", "clone https://user:s3cret@github.com/a/b failed", "clone https://github.com/a/b failed"},
		{"hides a GitHub token", "auth failed for ghp_abcdefghijklmnopqrstuvwxyz0123", "auth failed for [secret]"},
		{"hides a GitLab token", "bad glpat-abcdefghijklmnopqrst token", "bad [secret] token"},
		{"hides a bearer token", "sent Bearer abc.def-123", "sent Bearer [secret]"},
		{"hides key=value secrets", "failed password=hunter2 token: abc", "failed password=[secret] token: [secret]"},
		{"trims whitespace", "  padded \n", "padded"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SanitizeFailureReason(tc.in); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSanitizeFailureReasonFitsTheColumnWithoutSplittingCharacters(t *testing.T) {
	got := SanitizeFailureReason(strings.Repeat("é", FailureReasonMaxLen+50))
	if n := utf8.RuneCountInString(got); n != FailureReasonMaxLen {
		t.Errorf("got %d characters, want %d", n, FailureReasonMaxLen)
	}
	if !utf8.ValidString(got) {
		t.Error("truncation produced invalid UTF-8")
	}
}
