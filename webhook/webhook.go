package webhook

import (
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v01"
)

type Parser interface {
	Parse(r *http.Request) (*cloudevents.Event, error)
}
