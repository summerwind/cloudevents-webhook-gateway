package clair

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/summerwind/cloudevents-webhook-gateway/cloudevents"
)

type Webhook struct {
	Notification WebhookNotification `json:"Notification"`
}

type WebhookNotification struct {
	Name string `json:"Name"`
}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(req *http.Request) (*cloudevents.Event, error) {
	var w Webhook

	if req.Body == nil {
		return nil, errors.New("empty payload")
	}

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	err := decoder.Decode(&w)
	if err != nil {
		return nil, err
	}

	source := fmt.Sprintf("/notifications/%s", w.Notification.Name)
	s, err := url.Parse(source)
	if err != nil {
		return nil, err
	}

	ce := &cloudevents.Event{
		Type:            "com.coreos.clair.notify",
		Source:          *s,
		DataContentType: "application/json",
	}

	return ce, nil
}
