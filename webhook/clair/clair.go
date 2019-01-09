package clair

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v02"
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

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	err := decoder.Decode(&w)
	if err != nil {
		return nil, err
	}

	t := time.Now()

	source := fmt.Sprintf("/notifications/%s", w.Notification.Name)
	s, err := url.Parse(source)
	if err != nil {
		return nil, err
	}

	ce := &cloudevents.Event{
		Time:        &t,
		Type:        "com.coreos.clair.notify",
		Source:      *s,
		ContentType: "application/json",
	}

	return ce, nil
}
