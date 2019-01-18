package alertmanager

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go/v02"
	"github.com/prometheus/alertmanager/notify"
)

const (
	eventType   = "io.prometheus.alertmanager.alert"
	contentType = "application/json"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(req *http.Request) (*cloudevents.Event, error) {
	var msg notify.WebhookMessage

	if req.Body == nil {
		return nil, errors.New("empty payload")
	}

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	err := decoder.Decode(&msg)
	if err != nil {
		return nil, err
	}

	s, err := url.Parse(msg.ExternalURL)
	if err != nil {
		return nil, err
	}

	ce := &cloudevents.Event{
		Type:        eventType,
		Source:      *s,
		ContentType: contentType,
	}

	return ce, nil
}
