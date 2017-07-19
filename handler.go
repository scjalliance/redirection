package redirection

import "net/http"

// NewHandler returns a redirection handler for the given set of mappers.
func NewHandler(mapper ...Mapper) http.Handler {
	return MapperSet(mapper)
}
