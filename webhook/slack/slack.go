package slack

import (
	"errors"
	"net/http"
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go/v02"
)

const (
	eventType   = "com.slack.slash_command"
	contentType = "application/x-www-form-urlencoded"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(req *http.Request) (*cloudevents.Event, error) {
	if req.Body == nil {
		return nil, errors.New("empty payload")
	}

	req.ParseForm()
	defer req.Body.Close()

	command := req.FormValue("command")
	if command == "" {
		return nil, errors.New("enpty command")
	}

	tid := req.FormValue("trigger_id")
	if tid == "" {
		return nil, errors.New("empty trigger ID")
	}

	s, err := url.Parse(command)
	if err != nil {
		return nil, err
	}

	ce := &cloudevents.Event{
		ID:          tid,
		Type:        eventType,
		Source:      *s,
		ContentType: contentType,
	}

	return ce, nil
}
