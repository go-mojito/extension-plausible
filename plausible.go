package plausible

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-mojito/mojito/pkg/router"
)

var (
	Domain = ""
	URL    = "https://plausible.io"
)

func init() {
	router.RegisterHandlerArgFactory[PlausibleContext](func(ctx router.Context, next router.HandlerFunc) reflect.Value {
		return reflect.ValueOf(NewPlausibleContext(ctx))
	})
}

// PlausibleEvent defines the structure of the payload sent to the
// Plausible API for any event
type PlausibleEvent struct {
	Domain  string            `json:"d"`
	Event   string            `json:"n"`
	URL     string            `json:"u"`
	Referer *string           `json:"r"`
	Width   int               `json:"w"`
	Payload map[string]string `json:"p,omitempty"`
}

func SubmitEvent(event PlausibleEvent) error {
	jsonBody, err := json.Marshal(event)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/event", URL), "text/plain", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

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
