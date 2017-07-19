package redirection

import (
	"context"
	"fmt"
	"net/url"
	"sort"
)

// HostWhitelist returns a host whitelist function for the given mapper. The
// whitelist function will return a non-nil error if the given host name is not
// whitelisted by the mapper.
//
// In order for a host to be whitelisted the mapper must produce a result for
// that host's root path that has a weight greater than or equal to zero.
//
// The signature of the returned function is compatible with the HostPolicy
// type defined in the x/crypto/acme/autocert package.
func HostWhitelist(mapper Mapper) func(context.Context, string) error {
	fn := func(ctx context.Context, host string) error {
		u := url.URL{
			Host: host,
			Path: "/",
		}
		results := mapper.Map(&u, 0)
		sort.Sort(results)
		if len(results) > 0 && results[0].Weight >= 0 {
			return nil
		}
		return fmt.Errorf("host \"%s\" is not whitelisted for redirection", host)
	}
	return fn
}
