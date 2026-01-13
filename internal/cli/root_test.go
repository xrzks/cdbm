package cli

import (
	"strings"
	"testing"
)

func TestShellQuote(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal path", "/tmp/normal", "'/tmp/normal'"},
		{"path with space", "/tmp/with space", "'/tmp/with space'"},
		{"path with single quote", "/tmp/with'doublequote", "'/tmp/with'\\''doublequote'"},
		{"path with semicolon", "/tmp/with;semicolon", "'/tmp/with;semicolon'"},
		{"path with backtick", "/tmp/with`backtick`", "'/tmp/with`backtick`'"},
		{"path with dollar sign", "/tmp/with$dollar", "'/tmp/with$dollar'"},
		{"path with double quote", "/tmp/with\"quote\"", "'/tmp/with\"quote\"'"},
		{"path with newline", "/tmp/with\nnewline", "'/tmp/with\nnewline'"},
		{"path with carriage return", "/tmp/with\rcarriage", "'/tmp/with\rcarriage'"},
		{"path with backslash", "/tmp/with\\backslash", "'/tmp/with\\backslash'"},
		{"path with percent", "/tmp/with%percent", "'/tmp/with%percent'"},
		{"path with exclamation", "/tmp/with!exclamation", "'/tmp/with!exclamation'"},
		{"path with ampersand", "/tmp/with&ampersand", "'/tmp/with&ampersand'"},
		{"path with pipe", "/tmp/with(pipe)", "'/tmp/with(pipe)'"},
		{"path with redirect < ", "/tmp/with<redirect", "'/tmp/with<redirect'"},
		{"path with redirect >", "/tmp/with>redirect", "'/tmp/with>redirect'"},
		{"path with question mark", "/tmp/with?question", "'/tmp/with?question'"},
		{"path with bracket", "/tmp/with[bracket]", "'/tmp/with[bracket]'"},
		{"path with brace", "/tmp/with{brace}", "'/tmp/with{brace}'"},
		{"path with tilde", "/tmp/with~tilde", "'/tmp/with~tilde'"},
		{"path with hash", "/tmp/with#hash", "'/tmp/with#hash'"},
		{"path with caret", "/tmp/with^caret", "'/tmp/with^caret'"},
		{"path with pipe char", "/tmp/with|pipe", "'/tmp/with|pipe'"},
		{"path with at sign", "/tmp/with@sign", "'/tmp/with@sign'"},
		{"path with equals", "/tmp/with=equals", "'/tmp/with=equals'"},
		{"path with backslash", "/tmp/with\\slash", "'/tmp/with\\slash'"},
		{"path with tab", "/tmp/with\ttab", "'/tmp/with\ttab'"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := shellQuote(tc.input)
			if result != tc.expected {
				t.Errorf("shellQuote(%q) = %q, want %q", tc.input, result, tc.expected)
			}

			// Verify it's actually safe by checking for special chars outside quotes
			// After quoting, special chars should be inside single quotes
			if !strings.HasPrefix(result, "'") || !strings.HasSuffix(result, "'") {
				t.Errorf("Result not quoted: %q", result)
			}
		})
	}
}

func TestShellQuoteInjectionSafety(t *testing.T) {
	// Test that quoted output cannot execute commands
	testCases := []struct {
		name  string
		input string
	}{
		{"command substitution", "/tmp; rm -rf /"},
		{"backtick execution", "/tmp`touch /tmp/pwned`"},
		{"variable expansion", "/tmp$HOME"},
		{"newline injection", "/tmp\nwhoami"},
		{"carriage return", "/tmp\rwhoami"},
		{"command chaining", "/tmp && whoami"},
		{"pipe injection", "/tmp | whoami"},
		{"background job", "/tmp & whoami"},
		{"subshell", "/tmp $(whoami)"},
		{"logical or", "/tmp || whoami"},
		{"command continuation", "/tmp; whoami"},
		{"logical and", "/tmp && ls"},
		{"output redirect", "/tmp > /tmp/pwned"},
		{"input redirect", "/tmp < /etc/passwd"},
		{"command subshell", "/tmp `whoami`"},
		{"arithmetic expansion", "/tmp $((1+1))"},
		{"parameter expansion", "/tmp ${HOME}"},
		{"command substitution modern", "/tmp $(id)"},
		{"pipe to command", "/tmp | cat"},
		{"nested commands", "/tmp; ls; pwd"},
		{"background with output", "/tmp & > /dev/null"},
		{"complex injection", "/tmp; export EVIL=$(cat ~/.ssh/id_rsa)"},
		{"with double quote", "/tmp\"; rm -rf /"},
		{"with backtick and semicolon", "/tmp`rm -rf /`; whoami"},
		{"with dollar and parens", "/tmp$(whoami)"},
		{"with curly braces", "/tmp${HOME}"},
		{"multiple semicolons", "/tmp;;; rm -rf /"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			quoted := shellQuote(tc.input)

			// Should start and end with single quote
			if !strings.HasPrefix(quoted, "'") || !strings.HasSuffix(quoted, "'") {
				t.Errorf("Not properly quoted: %q -> %q", tc.input, quoted)
			}

			// Inner content should have unescaped single quotes escaped as '\''
			// This ensures the shell interprets it as a literal string
			inner := quoted[1 : len(quoted)-1]
			// Check no unescaped single quotes remain
			if strings.Contains(inner, "'") && !strings.Contains(inner, "'\\''") {
				t.Errorf("Unescaped single quote in: %q", quoted)
			}

			// Verify that when shell parses this, it's a single argument
			// Count single quotes: should be even (opening and closing pairs)
			quoteCount := strings.Count(quoted, "'")
			if quoteCount%2 != 0 {
				t.Errorf("Uneven number of quotes in: %q", quoted)
			}
		})
	}
}

func TestShellQuoteEmptyString(t *testing.T) {
	result := shellQuote("")
	if result != "''" {
		t.Errorf("shellQuote(\"\") = %q, want ''", result)
	}
}

func TestShellQuoteMultipleSingleQuotes(t *testing.T) {
	input := "/tmp/path'with'multiple'quotes"
	result := shellQuote(input)

	// Should escape each single quote individually
	if !strings.Contains(result, "'\\''") {
		t.Errorf("Single quotes not properly escaped in: %q", result)
	}

	// Count occurrences - should have 3 escaped single quotes
	escapedCount := strings.Count(result, "'\\''")
	if escapedCount != 3 {
		t.Errorf("Expected 3 escaped quotes in %q, got %d", result, escapedCount)
	}
}

func TestShellQuoteUnicode(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"emoji", "/tmp/📁"},
		{"chinese", "/tmp/目录"},
		{"arabic", "/tmp/دليل"},
		{"cyrillic", "/tmp/каталог"},
		{"emoji with spaces", "/tmp/📁 test 📂"},
		{"mixed unicode", "/tmp/test-目录-📁"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := shellQuote(tc.input)

			// Should be properly quoted
			if !strings.HasPrefix(result, "'") || !strings.HasSuffix(result, "'") {
				t.Errorf("Not properly quoted: %q -> %q", tc.input, result)
			}

			// Should contain the original input
			if !strings.Contains(result, tc.input) {
				t.Errorf("Result %q doesn't contain original input %q", result, tc.input)
			}
		})
	}
}
