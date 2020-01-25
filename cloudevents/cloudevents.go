package cloudevents

import (
	"net/url"
	"time"
)

type Event struct {
	ID     string
	Source url.URL
	Type   string

	DataContentType string
	DataSchema      url.URL
	Subject         string
	Time            *time.Time
}
