package plausible

import (
	"fmt"
	"strings"

	"github.com/go-mojito/mojito"
)

// PlausibleContext provides request-based functions to interact with Plausible
type PlausibleContext interface {
	mojito.Context

	// PageView will send a pageview event to the Plausible API for tracking visits
	PageView(pageUrl ...string) error

	// Trigger will send a custom event to the Plausible API for tracking goals
	Trigger(event string, payload map[string]string) error
}

// plausibleContext provides the implementation for PlausibleContext
type plausibleContext struct {
	mojito.Context
}

// domain will resolve the domain the event is being created for
func (p *plausibleContext) IsDomain() bool {
	return strings.ToLower(p.Request().GetRequest().URL.Host) == strings.ToLower(domain)
}

// createEvent will create an empty event with request-based values filled in
func (p *plausibleContext) createEvent(eventName string) *PlausibleEvent {
	url := p.Request().GetRequest().URL
	referer := p.Request().GetRequest().Referer()
	forwardedFor := p.Request().GetRequest().Header.Get("X-Forwarded-For")
	realIP := p.Request().GetRequest().Header.Get("X-Real-IP")
	event := &PlausibleEvent{
		Domain:  strings.ToLower(p.Request().GetRequest().URL.Host),
		Event:   eventName,
		URL:     fmt.Sprintf("http://%s%s", strings.ToLower(p.Request().GetRequest().URL.Host), url.Path),
		Width:   guessWidthFromUA(p.Request().GetRequest().UserAgent()),
		Payload: nil,
	}
	if referer != "" {
		event.Referer = &referer
	}
	if realIP != "" {
		event.IP = &realIP
	}
	if forwardedFor != "" {
		event.IP = &forwardedFor
	}
	return event
}

// PageView will send a pageview event to the Plausible API for tracking visits
func (p *plausibleContext) PageView(pageUrl ...string) error {
	event := p.createEvent("pageview")
	if pageUrl != nil && len(pageUrl) > 0 {
		event.URL = fmt.Sprintf("http://%s/%s", strings.ToLower(p.Request().GetRequest().URL.Host), strings.Join(pageUrl, "/"))
	}
	return SubmitEvent(*event)
}

// Trigger will send a custom event to the Plausible API for tracking goals
func (p *plausibleContext) Trigger(eventName string, payload map[string]string) error {
	event := p.createEvent(eventName)
	event.Payload = payload
	return SubmitEvent(*event)
}

// NewPlausibleContext will create a new instance of the default implementation for PlausibleContext
func NewPlausibleContext(ctx mojito.Context) PlausibleContext {
	return &plausibleContext{
		Context: ctx,
	}
}
