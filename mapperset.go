package redirection

import (
	"log"
	"net/http"
	"net/url"
	"sort"
)

// MapperSet is a set of redirection mappers.
type MapperSet []Mapper

// Map returns an unsorted set of redirection results for the given url, with
// the provided weight added to each result.
//
// The results from each mapper in the set will be included.
func (set MapperSet) Map(u *url.URL, weight int) (results ResultSet) {
	for _, mapper := range set {
		if subresults := mapper.Map(u, weight); len(subresults) > 0 {
			results = append(results, subresults...)
		}
	}
	return
}

func (set MapperSet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := urlFromRequest(r)
	results := set.Map(u, 0)
	sort.Stable(results)

	// TODO: Move this output to some sort of configurable logging target
	log.Printf("\"%s\" %+v", u, results)

	if len(results) == 0 {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	result := results[0]

	switch result.Code {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect, http.StatusPermanentRedirect:
		http.Redirect(w, r, result.URL, result.Code)
	default:
		http.Error(w, http.StatusText(result.Code), result.Code)
	}
}

func urlFromRequest(r *http.Request) *url.URL {
	u := *r.URL
	u.Host = r.Host
	return &u
}
