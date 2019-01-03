package github

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	Secret      = "test"
	EventID     = "72d3162e-cc78-11e3-81ab-4c9367dc0958"
	ContentType = "application/json"
)

func loadFixture(name string) ([]byte, error) {
	_, fn, _, _ := runtime.Caller(0)
	fx := filepath.Join(filepath.Dir(fn), "fixtures", fmt.Sprintf("%s.json", name))
	return ioutil.ReadFile(fx)
}

func getSignature(payload, secret []byte) string {
	mac := hmac.New(sha1.New, secret)
	mac.Write(payload)
	return fmt.Sprintf("sha1=%s", hex.EncodeToString(mac.Sum(nil)))
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
	req.Header.Set("X-GitHub-Event", name)
	req.Header.Set("X-GitHub-Delivery", EventID)
	req.Header.Set("X-Hub-Signature", getSignature(body, []byte(Secret)))

	return req, nil
}

func TestParse(t *testing.T) {
	tests := []struct {
		eventType string
		ceType    string
		ceSource  string
	}{
		{"check_run", "com.github.check_run", "https://api.github.com/repos/github/hello-world/check-runs/4"},
		{"check_suite", "com.github.check_suite", "https://api.github.com/repos/github/hello-world/check-suites/5"},
		{"commit_comment", "com.github.commit_comment", "https://api.github.com/repos/Codertocat/Hello-World/comments/29186860"},
		{"create", "com.github.create", "https://api.github.com/repos/Codertocat/Hello-World"},
		{"delete", "com.github.delete", "https://api.github.com/repos/Codertocat/Hello-World"},
		{"deployment", "com.github.deployment", "https://api.github.com/repos/Codertocat/Hello-World/deployments/87972451"},
		{"deployment_status", "com.github.deployment_status", "https://api.github.com/repos/Codertocat/Hello-World/deployments/87972451"},
		{"fork", "com.github.fork", "https://api.github.com/repos/Octocoders/Hello-World"},
		{"gollum", "com.github.gollum", "https://api.github.com/repos/Codertocat/Hello-World"},
		{"installation", "com.github.installation", "https://github.com/settings/installations/2"},
		{"installation_repositories", "com.github.installation_repositories", "https://github.com/settings/installations/2"},
		{"issue_comment", "com.github.issue_comment", "https://api.github.com/repos/Codertocat/Hello-World/issues/comments/393304133"},
		{"issues", "com.github.issues", "https://api.github.com/repos/Codertocat/Hello-World/issues/2"},
		{"label", "com.github.label", "https://api.github.com/repos/Codertocat/Hello-World/labels/:bug:%20Bugfix"},
		{"marketplace_purchase", "com.github.marketplace_purchase", "https://api.github.com/users/username"},
		{"member", "com.github.member", "https://api.github.com/users/octocat"},
		{"membership", "com.github.membership", "https://api.github.com/teams/2723476"},
		{"milestone", "com.github.milestone", "https://api.github.com/repos/Codertocat/Hello-World/milestones/1"},
		{"organization", "com.github.organization", "https://api.github.com/orgs/Octocoders"},
		{"org_block", "com.github.org_block", "https://api.github.com/orgs/Octocoders"},
		{"page_build", "com.github.page_build", "https://api.github.com/repos/Codertocat/Hello-World/pages/builds/91762186"},
		{"project_card", "com.github.project_card", "https://api.github.com/projects/columns/cards/10189042"},
		{"project_column", "com.github.project_column", "https://api.github.com/projects/columns/2803722"},
		{"project", "com.github.project", "https://api.github.com/projects/1547122"},
		{"public", "com.github.public", "https://api.github.com/repos/Codertocat/Hello-World"},
		{"pull_request", "com.github.pull_request", "https://api.github.com/repos/Codertocat/Hello-World/pulls/1"},
		{"pull_request_review", "com.github.pull_request_review", "https://api.github.com/repos/Codertocat/Hello-World/pulls/1"},
		{"pull_request_review_comment", "com.github.pull_request_review_comment", "https://api.github.com/repos/Codertocat/Hello-World/pulls/comments/191908831"},
		{"push", "com.github.push", "https://api.github.com/repos/Codertocat/Hello-World/git/refs/tags/simple-tag"},
		{"release", "com.github.release", "https://api.github.com/repos/Codertocat/Hello-World/releases/11248810"},
		{"repository", "com.github.repository", "https://api.github.com/repos/Codertocat/Hello-World"},
		{"status", "com.github.status", "https://api.github.com/repos/Codertocat/Hello-World/commits/a10867b14bb761a232cd80139fbd4c0d33264240"},
		{"team", "com.github.team", "https://api.github.com/teams/2723476"},
		{"team_add", "com.github.team_add", "https://api.github.com/teams/2723476"},
		{"watch", "com.github.watch", "https://api.github.com/repos/Codertocat/Hello-World"},
	}

	for _, test := range tests {
		req, err := newRequest(test.eventType)
		if err != nil {
			t.Fatalf("[%s] invalid request: %v", test.eventType, err)
		}

		p := NewParser(Secret)
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
