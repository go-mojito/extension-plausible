package plausible

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: time.Second * 10,
}

// Event defines the structure of the payload sent to the
// Plausible API for any event
type Event struct {
	Domain  string            `json:"d"`
	Event   string            `json:"n"`
	IP      *string           `json:"-"`
	URL     string            `json:"u"`
	Referer *string           `json:"r"`
	Width   int               `json:"w"`
	Payload map[string]string `json:"p,omitempty"`
}

// SubmitEvent will POST an event to the plausible API
func SubmitEvent(event Event) error {
	jsonBody, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/event", url), bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "text/plain")
	if event.IP != nil {
		req.Header.Set("X-Forwarded-For", *event.IP)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	body := string(bodyBytes)
	if !strings.Contains(body, "ok") {
		return errors.New(body)
	}
	return nil
}
