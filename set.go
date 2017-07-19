package redirection

import "net/url"

// Set is a set of redirection elements.
type Set []Element

// Map returns an unsorted set of redirection results for the given url, with
// the provided weight added to each result.
//
// The results from each element in the set will be included.
func (set Set) Map(u *url.URL, weight int) (results ResultSet) {
	for i := range set {
		if subresults := set[i].Map(u, weight); len(subresults) > 0 {
			results = append(results, subresults...)
		}
	}
	return
}
