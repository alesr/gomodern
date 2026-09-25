// Package legacy contains deliberately legacy Go idioms used as fixtures for
// gomodern integration tests. Each construct is behavior-preserving and is
// modernized by the modernize analyzer suite.
package legacy

import "sort"

// TakesInterface accepts the empty interface type (pre-Go 1.18 "any" alias).
func TakesInterface(v interface{}) interface{} {
	return v
}

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

// SortInts sorts a slice using the pre-1.21 sort.Slice idiom.
func SortInts(s []int) {
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
}
