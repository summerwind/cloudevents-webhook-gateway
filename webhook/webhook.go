package webhook

import (
	"net/http"

	"github.com/summerwind/cloudevents-webhook-gateway/cloudevents"
)

type Parser interface {
	Parse(r *http.Request) (*cloudevents.Event, error)
}
