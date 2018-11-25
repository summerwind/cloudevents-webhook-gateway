package alertmanager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v01"
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

	ce := &cloudevents.Event{
		EventTime:        &t,
		EventType:        fmt.Sprintf("io.prometheus.alertmanager.%s", msg.Status),
		EventTypeVersion: msg.Version,
		Source:           fmt.Sprintf("/groupKey/%s", msg.GroupKey),
		ContentType:      "application/json",
	}

	return ce, nil
}
