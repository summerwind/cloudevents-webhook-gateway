package proxy

import (
	"errors"
	"net/http"
)

type Transport struct {
	base http.RoundTripper
}

func NewTransport() Transport {
	return Transport{base: http.DefaultTransport}
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("CE-ID") == "" {
		return nil, errors.New("invalid request")
	}
	return t.base.RoundTrip(req)
}
