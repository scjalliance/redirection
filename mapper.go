package redirection

import "net/url"

// Mapper is a redirection mapper. It returns an unsorted set of redirection
// results for the given url, with the provided weight added to each result.
type Mapper interface {
	Map(u *url.URL, weight int) (results ResultSet)
}
