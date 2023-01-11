package plausible

import (
	"reflect"
	"regexp"

	"github.com/go-mojito/mojito/pkg/router"
)

var (
	domain              = ""
	enforceDomainFilter = true
	filters             = []regexp.Regexp{}
	url                 = "https://plausible.io"
)

func init() {
	router.RegisterHandlerArgFactory[Context](func(ctx router.Context, next router.HandlerFunc) (reflect.Value, error) {
		return reflect.ValueOf(NewContext(ctx)), nil
	})
}

// Configure sets the domain of the plausible site that will be reported for.
// This is required to enable the extension
func Configure(siteDomain string) {
	domain = siteDomain
}

// EnforceDomain enables or disables the filter that will prevent tracking from happening
// when the request host does not match the configured domain
func EnforceDomain(enforce bool) {
	enforceDomainFilter = enforce
}

// Ignore will skip matching URL paths when using the automatic middleware tracking.
// This is useful if you want to filter out assets, admin areas, etc from your stats.
func Ignore(r *regexp.Regexp) {
	filters = append(filters, *r)
}

// SetInstanceURL will change the API base URL to the given URL.
// This is useful if you self-host plausible and you want to use that instance
func SetInstanceURL(instanceURL string) {
	url = instanceURL
}
