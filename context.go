package plausible

import (
	"fmt"
	"strings"

	"github.com/go-mojito/mojito"
)

type PlausibleContext interface {
	mojito.Context

	// PageView will send a pageview event to the Plausible API for tracking visits
	PageView(pageUrl ...string) error

	// Trigger will send a custom event to the Plausible API for tracking goals
	Trigger(event string, payload map[string]string) error
}

type plausibleContext struct {
	mojito.Context
}

// domain will resolve the domain the event is being created for
func (p *plausibleContext) domain() string {
	if Domain != "" {
		return Domain
	}
	return p.Request().GetRequest().URL.Host
}

// createEvent will create an empty event with request-based values filled in
func (p *plausibleContext) createEvent(eventName string) *PlausibleEvent {
	url := p.Request().GetRequest().URL
	referer := p.Request().GetRequest().Referer()
	event := &PlausibleEvent{
		Domain:  p.domain(),
		Event:   eventName,
		URL:     fmt.Sprintf("http://%s%s", p.domain(), url.Path),
		Width:   guessWidthFromUA(p.Request().GetRequest().UserAgent()),
		Payload: nil,
	}
	if referer != "" {
		event.Referer = &referer
	}
	return event
}

// PageView will send a pageview event to the Plausible API for tracking visits
func (p *plausibleContext) PageView(pageUrl ...string) error {
	event := p.createEvent("pageview")
	if pageUrl != nil && len(pageUrl) > 0 {
		event.URL = fmt.Sprintf("http://%s/%s", p.domain(), strings.Join(pageUrl, "/"))
	}
	return SubmitEvent(*event)
}

// Trigger will send a custom event to the Plausible API for tracking goals
func (p *plausibleContext) Trigger(eventName string, payload map[string]string) error {
	event := p.createEvent(eventName)
	event.Payload = payload
	return SubmitEvent(*event)
}

func NewPlausibleContext(ctx mojito.Context) PlausibleContext {
	return &plausibleContext{
		Context: ctx,
	}
}
