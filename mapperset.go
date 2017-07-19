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
	results := set.Map(r.URL, 0)
	if len(results) == 0 {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	sort.Stable(results)

	result := results[0]
	log.Printf("\"%s\" %+v", r.URL, results)

	switch result.Code {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect, http.StatusPermanentRedirect:
		http.Redirect(w, r, result.URL, result.Code)
	default:
		http.Error(w, http.StatusText(result.Code), result.Code)
	}
}
