package redirection

import "fmt"

// Result is a weighted target that has been successfully mapped.
type Result struct {
	Target
	Weight int `json:"weight"`
}

// String encodes result as a string.
func (r Result) String() string {
	return fmt.Sprintf("{%d %s %d}", r.Code, r.URL, r.Weight)
}

// Supersedes returns true if r supersedes c.
func (r Result) Supersedes(c Result) bool {
	if r.Weight > c.Weight {
		return true
	}
	return false
}
