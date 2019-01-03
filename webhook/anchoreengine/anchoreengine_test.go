package anchoreengine

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
	ContentType = "application/json"
)

func loadFixture(name string) ([]byte, error) {
	_, fn, _, _ := runtime.Caller(0)
	fx := filepath.Join(filepath.Dir(fn), "fixtures", fmt.Sprintf("%s.json", name))
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
	tests := []struct {
		eventType string
		ceType    string
		ceSource  string
	}{
		{"analysis_update", "com.anchore.anchore-engine.analysis_update", "/v1/subscriptions?subscription_key=docker.io/dnurmi/testrepo:latest"},
		{"tag_update", "com.anchore.anchore-engine.tag_update", "/v1/subscriptions?subscription_key=docker.io/dnurmi/testrepo:latest"},
		{"policy_eval", "com.anchore.anchore-engine.policy_eval", "/v1/subscriptions?subscription_key=docker.io/dnurmi/testrepo:latest"},
		{"vuln_update", "com.anchore.anchore-engine.vuln_update", "/v1/subscriptions?subscription_key=docker.io/dnurmi/testrepo:latest"},
	}

	for _, test := range tests {
		req, err := newRequest(test.eventType)
		if err != nil {
			t.Fatalf("[%s] invalid request: %v", test.eventType, err)
		}

		p := NewParser()
		ce, err := p.Parse(req)
		if err != nil {
			t.Fatalf("[%s] parser error: %v", test.eventType, err)
		}

		if ce.Type != test.ceType {
			t.Errorf("[%s] invalid type: %v", test.eventType, ce.Type)
		}
		if ce.Source.String() != test.ceSource {
			t.Errorf("[%s] invalid source: %v", test.eventType, ce.Source)
		}
	}
}
