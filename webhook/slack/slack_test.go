package slack

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	ContentType = "application/x-www-form-urlencoded"
)

func loadFixture(name string) ([]byte, error) {
	_, fn, _, _ := runtime.Caller(0)
	fx := filepath.Join(filepath.Dir(fn), "fixtures", fmt.Sprintf("%s.txt", name))
	return ioutil.ReadFile(fx)
}

func newRequest(name string) (*http.Request, error) {
	body, err := loadFixture(name)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("Content-Length", string(len(body)))

	return req, nil
}

func TestParse(t *testing.T) {
	req, err := newRequest("slash_command")
	if err != nil {
		t.Fatalf("invalid request: %v", err)
	}

	p := NewParser()
	ce, err := p.Parse(req)
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	if ce.Type != "com.slack.slash_command" {
		t.Errorf("invalid type: %v", ce.Type)
	}
	if ce.Source.String() != "https://example.slack.com/messages/C2147483705" {
		t.Errorf("invalid source: %v", ce.Source)
	}
}
