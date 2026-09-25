// Package legacy120 contains legacy idioms that require Go 1.21+ to
// modernize. With a `go 1.20` directive, the minmax and slicescontains
// analyzers must be suppressed, so a fix run must leave this file unchanged.
package legacy120

// MaxInt implements max by hand (Go 1.21 introduced the builtin max).
func MaxInt(a, b int) int {
	var m int
	if a > b {
		m = a
	} else {
		m = b
	}
	return m
}

// Contains checks slice membership with a manual loop (Go 1.21 slices.Contains).
func Contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
