package pattern

// Set is a set of patterns.
type Set []Pattern

// Match returns true if one or more patterns in the set match the string.
//
// If one or more of the matching patterns are a regular expression with
// capturing groups, the combined set of all the non-empty submatches will be
// returned as well.
func (set Set) Match(s string) (submatches []string, matched bool) {
	seen := make(map[string]bool)
	for i := range set {
		if subm, ok := set[i].Match(s); ok {
			matched = true
			for _, submatch := range subm {
				if !seen[submatch] {
					seen[submatch] = true
					submatches = append(submatches, submatch)
				}
			}
		}
	}
	return
}
