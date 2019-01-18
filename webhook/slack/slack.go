package slack

import (
	"errors"
	"fmt"
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
	req.ParseForm()
	defer req.Body.Close()

	td := req.FormValue("team_domain")
	if td == "" {
		return nil, errors.New("Empty team domain")
	}

	cid := req.FormValue("channel_id")
	if cid == "" {
		return nil, errors.New("Empty channel ID")
	}

	tid := req.FormValue("trigger_id")
	if tid == "" {
		return nil, errors.New("Empty trigger ID")
	}

	s, err := url.Parse(fmt.Sprintf("https://%s.slack.com/messages/%s", td, cid))
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
