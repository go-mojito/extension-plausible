package plausible

// PlausibleEvent defines the structure of the payload sent to the
// Plausible API for any event
type PlausibleEvent struct {
	Domain  string            `json:"d"`
	Event   string            `json:"n"`
	IP      *string           `json:"-"`
	URL     string            `json:"u"`
	Referer *string           `json:"r"`
	Width   int               `json:"w"`
	Payload map[string]string `json:"p,omitempty"`
}
