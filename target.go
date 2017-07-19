package redirection

import "net/http"

// Target holds the target URL and HTTP status code for a redirect.
//
// If the URL is empty no redirection will occur, but the HTTP status code
// will be returned.
type Target struct {
	URL  string `json:"url"`
	Code int    `json:"code"`
}

// Valid reports whether t has a valid status code.
func (t *Target) Valid() bool {
	if t.Code == 0 {
		return false
	}
	return http.StatusText(t.Code) != ""
}
