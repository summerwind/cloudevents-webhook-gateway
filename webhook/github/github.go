package github

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v02"
	"github.com/google/go-github/v21/github"
)

type Parser struct {
	secret []byte
}

func NewParser(secret string) *Parser {
	return &Parser{
		secret: []byte(secret),
	}
}

func (p *Parser) Parse(req *http.Request) (*cloudevents.Event, error) {
	if req.Body == nil {
		return nil, errors.New("empty payload")
	}

	payload, err := github.ValidatePayload(req, p.secret)
	if err != nil {
		return nil, err
	}

	webHookType := github.WebHookType(req)
	event, err := github.ParseWebHook(webHookType, payload)
	if err != nil {
		return nil, err
	}

	var (
		action    string
		source    string
		eventType string
	)

	switch event := event.(type) {
	case *github.CheckRunEvent:
		source = event.CheckRun.GetURL()
		action = event.GetAction()
	case *github.CheckSuiteEvent:
		source = event.CheckSuite.GetURL()
	case *github.CommitCommentEvent:
		source = event.Comment.GetURL()
		action = event.GetAction()
	case *github.CreateEvent:
		source = event.Repo.GetURL()
	case *github.DeleteEvent:
		source = event.Repo.GetURL()
	case *github.DeploymentEvent:
		source = event.Deployment.GetURL()
	case *github.DeploymentStatusEvent:
		source = event.Deployment.GetURL()
	case *github.ForkEvent:
		source = event.Forkee.GetURL()
	case *github.GollumEvent:
		source = event.Repo.GetURL()
	case *github.InstallationEvent:
		source = event.Installation.GetHTMLURL()
		action = event.GetAction()
	case *github.InstallationRepositoriesEvent:
		source = event.Installation.GetHTMLURL()
		action = event.GetAction()
	case *github.IssueCommentEvent:
		source = event.Comment.GetURL()
		action = event.GetAction()
	case *github.IssuesEvent:
		source = event.Issue.GetURL()
		action = event.GetAction()
	case *github.LabelEvent:
		source = event.Label.GetURL()
		action = event.GetAction()
	case *github.MarketplacePurchaseEvent:
		source = event.Sender.GetURL()
		action = event.GetAction()
	case *github.MemberEvent:
		source = event.Member.GetURL()
		action = event.GetAction()
	case *github.MembershipEvent:
		source = event.Team.GetURL()
		action = event.GetAction()
	case *github.MilestoneEvent:
		source = event.Milestone.GetURL()
		action = event.GetAction()
	case *github.OrganizationEvent:
		source = event.Organization.GetURL()
		action = event.GetAction()
	case *github.OrgBlockEvent:
		source = event.Organization.GetURL()
		action = event.GetAction()
	case *github.PageBuildEvent:
		source = event.Build.GetURL()
	case *github.ProjectCardEvent:
		source = event.ProjectCard.GetURL()
		action = event.GetAction()
	case *github.ProjectColumnEvent:
		source = event.ProjectColumn.GetURL()
		action = event.GetAction()
	case *github.ProjectEvent:
		source = event.Project.GetURL()
		action = event.GetAction()
	case *github.PublicEvent:
		source = event.Repo.GetURL()
	case *github.PullRequestReviewCommentEvent:
		source = event.Comment.GetURL()
		action = event.GetAction()
	case *github.PullRequestReviewEvent:
		source = event.PullRequest.GetURL()
		action = event.GetAction()
	case *github.PullRequestEvent:
		source = event.PullRequest.GetURL()
		action = event.GetAction()
	case *github.PushEvent:
		// API URL is not set in "repository.url", need to generate URL from statuses URL.
		base, err := url.Parse(event.Repo.GetStatusesURL())
		if err != nil {
			return nil, err
		}
		ref, err := url.Parse(fmt.Sprintf("../git/%s", event.GetRef()))
		if err != nil {
			return nil, err
		}
		source = base.ResolveReference(ref).String()
	case *github.RepositoryEvent:
		source = event.Repo.GetURL()
		action = event.GetAction()
	case *github.ReleaseEvent:
		source = event.Release.GetURL()
		action = event.GetAction()
	case *github.StatusEvent:
		source = event.Commit.GetURL()
	case *github.TeamEvent:
		source = event.Team.GetURL()
		action = event.GetAction()
	case *github.TeamAddEvent:
		source = event.Team.GetURL()
	case *github.WatchEvent:
		source = event.Repo.GetURL()
		action = event.GetAction()
	}

	if source == "" {
		return nil, errors.New("unsupported event type")
	}

	if action == "" {
		eventType = fmt.Sprintf("com.github.%s", webHookType)
	} else {
		eventType = fmt.Sprintf("com.github.%s.%s", webHookType, action)
	}

	t := time.Now()

	s, err := url.Parse(source)
	if err != nil {
		return nil, err
	}

	ce := &cloudevents.Event{
		ID:          github.DeliveryID(req),
		Time:        &t,
		Type:        eventType,
		Source:      *s,
		ContentType: "application/json",
	}

	return ce, nil
}
