package redirection

import (
	"net/url"

	"github.com/scjalliance/redirection/pattern"
)

// Element is an element of redirection logic. It maps URLs to redirection
// targets, is capable of matching hosts and/or paths, and can be formed into a
// processing tree.
type Element struct {
	Host     pattern.Set `json:"host"`   // Must match URL host if not empty
	Path     pattern.Set `json:"path"`   // Must match URL path if not empty
	Mapped   Map         `json:"map"`    // Set of keyed sub-elements (uses pattern subgroup match)
	Elements Set         `json:"match"`  // Set of sub-elements
	Weight   int         `json:"weight"` // Weight applied to all targets within this element
	Target
}

// Map returns an unsorted set of redirection results for the given url, with
// the provided weight added to each result.
//
// The results of any matching sub-elements will be included.
func (e *Element) Map(u *url.URL, weight int) (results ResultSet) {
	weight += e.Weight

	var keys []string

	if len(e.Host) > 0 {
		hostKeys, ok := e.Host.Match(u.Hostname())
		if !ok {
			return nil
		}
		keys = append(keys, hostKeys...)
	}

	if len(e.Path) > 0 {
		pathKeys, ok := e.Path.Match(u.Path)
		if !ok {
			return nil
		}
		keys = append(keys, pathKeys...)
	}

	if len(e.Mapped) > 0 {
		elements := e.Mapped.Match(keys)
		if subresults := elements.Map(u, weight); len(subresults) > 0 {
			results = append(results, subresults...)
		}
	}

	if len(e.Elements) > 0 {
		if subresults := e.Elements.Map(u, weight); len(subresults) > 0 {
			results = append(results, subresults...)
		}
	}

	if e.Target.Valid() {
		results = append(results, Result{
			Target: e.Target,
			Weight: weight,
		})
	}

	return
}
