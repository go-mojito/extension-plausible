package plausible

import (
	"github.com/go-mojito/mojito/pkg/router"
)

// Context provides request-based functions to interact with Plausible
type Context interface {
	router.Context

	// PageView will send a pageview event to the Plausible API for tracking visits
	PageView(pageUrl ...string) error

	// Trigger will send a custom event to the Plausible API for tracking goals
	Trigger(event string, payload map[string]string) error
}

// NewContext will create a new instance of the default implementation for PlausibleContext
func NewContext(ctx router.Context) Context {
	return &context{
		Context: ctx,
	}
}
