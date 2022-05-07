package plausible

import (
	"fmt"
	"strings"

	"github.com/go-mojito/mojito"
)

// context provides the implementation for PlausibleContext
type context struct {
	mojito.Context
}

// domain will resolve the domain the event is being created for
func (p *context) IsDomain() bool {
	return strings.ToLower(p.Request().GetRequest().Host) == strings.ToLower(domain)
}

// createEvent will create an empty event with request-based values filled in
func (p *context) createEvent(eventName string) *Event {
	url := p.Request().GetRequest().URL
	referer := p.Request().GetRequest().Referer()
	forwardedFor := p.Request().GetRequest().Header.Get("X-Forwarded-For")
	realIP := p.Request().GetRequest().Header.Get("X-Real-IP")
	event := &Event{
		Domain:  domain,
		Event:   eventName,
		URL:     fmt.Sprintf("http://%s%s", domain, url.Path),
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
func (p *context) PageView(pageUrl ...string) error {
	if enforceDomainFilter && !p.IsDomain() {
		return nil
	}
	event := p.createEvent("pageview")
	if pageUrl != nil && len(pageUrl) > 0 {
		event.URL = fmt.Sprintf("http://%s/%s", domain, strings.Join(pageUrl, "/"))
	}
	return SubmitEvent(*event)
}

// Trigger will send a custom event to the Plausible API for tracking goals
func (p *context) Trigger(eventName string, payload map[string]string) error {
	if enforceDomainFilter && !p.IsDomain() {
		return nil
	}
	event := p.createEvent(eventName)
	event.Payload = payload
	return SubmitEvent(*event)
}
