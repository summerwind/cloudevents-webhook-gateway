package dockerhub

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v02"
)

type Webhook struct {
	Repository WebhookRepository `json:"repository"`
}

type WebhookRepository struct {
	RepoURL string `json:"repo_url"`
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

	s, err := url.Parse(w.Repository.RepoURL)
	if err != nil {
		return nil, err
	}

	ce := &cloudevents.Event{
		Time:        &t,
		Type:        "com.docker.hub.push",
		Source:      *s,
		ContentType: "application/json",
	}

	return ce, nil
}
