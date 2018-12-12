package alertmanager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v02"
	"github.com/prometheus/alertmanager/notify"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(req *http.Request) (*cloudevents.Event, error) {
	var msg notify.WebhookMessage

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	err := decoder.Decode(&msg)
	if err != nil {
		return nil, err
	}

	t := time.Now()

	s, err := url.Parse(msg.ExternalURL)
	if err != nil {
		return nil, err
	}

	ce := &cloudevents.Event{
		Time:        &t,
		Type:        fmt.Sprintf("io.prometheus.alertmanager.%s", msg.Status),
		Source:      *s,
		ContentType: "application/json",
	}

	return ce, nil
}
