package anchoreengine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go/v02"
)

type Webhook struct {
	Data WebhookData `json:"data"`
}

type WebhookData struct {
	NotificationType    string                     `json:"notification_type"`
	NotificationPayload WebhookNotificationPayload `json:"notification_payload"`
}

type WebhookNotificationPayload struct {
	NotificationID  string `json:"notificationId"`
	SubscriptionKey string `json:"subscription_key"`
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

	source := fmt.Sprintf("/v1/subscriptions?subscription_key=%s", w.Data.NotificationPayload.SubscriptionKey)
	s, err := url.Parse(source)
	if err != nil {
		return nil, err
	}

	ce := &cloudevents.Event{
		ID:          w.Data.NotificationPayload.NotificationID,
		Type:        fmt.Sprintf("com.anchore.anchore-engine.%s", w.Data.NotificationType),
		Source:      *s,
		ContentType: "application/json",
	}

	return ce, nil
}
