package redirection

// Map is a map of redirection elements.
type Map map[string]Element

// Match returns the set of elements that match one of the given keys.
func (m Map) Match(keys []string) (elements Set) {
	for _, key := range keys {
		if element, ok := m[key]; ok {
			elements = append(elements, element)
		}
	}
	return
}
