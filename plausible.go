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
	"time"

	"github.com/go-mojito/mojito/pkg/router"
)

var (
	client = &http.Client{
		Timeout: time.Second * 10,
	}
	domain              = ""
	enforceDomainFilter = true
	url                 = "https://plausible.io"
)

func init() {
	router.RegisterHandlerArgFactory[PlausibleContext](func(ctx router.Context, next router.HandlerFunc) reflect.Value {
		return reflect.ValueOf(NewPlausibleContext(ctx))
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

// SetInstanceURL will change the API base URL to the given URL.
// This is useful if you self-host plausible and you want to use that instance
func SetInstanceURL(instanceURL string) {
	url = instanceURL
}

// SubmitEvent will POST an event to the plausible API
func SubmitEvent(event PlausibleEvent) error {
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
