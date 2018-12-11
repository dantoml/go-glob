package glob

import (
	"strings"
	"unicode/utf8"
)

// The character which is treated like a glob
const GLOB = "*"

// Glob will test a string pattern, potentially containing globs,  against a
// subject string. The result is a true/false, determining whether or not the
// glob pattern matched the subject text
func Glob(pattern, subj string) bool {
	match, _ := GlobWithDifference(pattern, subj)
	return match
}

// GlobWithDifference will test a string pattern, potentially containing globs,
// against a subject string. The result is a true/false, determining
// whether or not the glob pattern matched the subject text, and the count of
// characters matched by the glob.
func GlobWithDifference(pattern, subj string) (bool, int) {
	// Empty pattern can only match empty subject
	if pattern == "" {
		return subj == pattern, 0
	}

	// If the pattern _is_ a glob, it matches everything
	if pattern == GLOB {
		return true, utf8.RuneCountInString(subj)
	}

	parts := strings.Split(pattern, GLOB)

	if len(parts) == 1 {
		// No globs in pattern, so test for equality
		return subj == pattern, 0
	}

	difference := utf8.RuneCountInString(subj)
	leadingGlob := strings.HasPrefix(pattern, GLOB)
	trailingGlob := strings.HasSuffix(pattern, GLOB)
	end := len(parts) - 1

	// Go over the leading parts and ensure they match.
	for i := 0; i < end; i++ {
		idx := strings.Index(subj, parts[i])

		switch i {
		case 0:
			// Check the first section. Requires special handling.
			if !leadingGlob && idx != 0 {
				return false, 0
			}
		default:
			// Check that the middle parts match.
			if idx < 0 {
				return false, 0
			}
		}

		difference -= utf8.RuneCountInString(parts[i])
		// Trim evaluated text from subj as we loop over the pattern.
		subj = subj[idx+len(parts[i]):]
	}

	// Reached the last section. Requires special handling.
	match := false
	if trailingGlob {
		match = true
	} else if strings.HasSuffix(subj, parts[end]) {
		difference -= utf8.RuneCountInString(parts[end])
		match = true
	}

	return match, difference
}
