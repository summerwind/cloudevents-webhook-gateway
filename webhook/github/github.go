package github

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-github/v21/github"
	"github.com/summerwind/cloudevents-webhook-gateway/cloudevents"
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
	var source string

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

	switch event := event.(type) {
	case *github.CheckRunEvent:
		source = event.CheckRun.GetURL()
	case *github.CheckSuiteEvent:
		source = event.CheckSuite.GetURL()
	case *github.CommitCommentEvent:
		source = event.Comment.GetURL()
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
	case *github.InstallationRepositoriesEvent:
		source = event.Installation.GetHTMLURL()
	case *github.IssueCommentEvent:
		source = event.Comment.GetURL()
	case *github.IssuesEvent:
		source = event.Issue.GetURL()
	case *github.LabelEvent:
		source = event.Label.GetURL()
	case *github.MarketplacePurchaseEvent:
		source = event.Sender.GetURL()
	case *github.MemberEvent:
		source = event.Member.GetURL()
	case *github.MembershipEvent:
		source = event.Team.GetURL()
	case *github.MilestoneEvent:
		source = event.Milestone.GetURL()
	case *github.OrganizationEvent:
		source = event.Organization.GetURL()
	case *github.OrgBlockEvent:
		source = event.Organization.GetURL()
	case *github.PageBuildEvent:
		source = event.Build.GetURL()
	case *github.ProjectCardEvent:
		source = event.ProjectCard.GetURL()
	case *github.ProjectColumnEvent:
		source = event.ProjectColumn.GetURL()
	case *github.ProjectEvent:
		source = event.Project.GetURL()
	case *github.PublicEvent:
		source = event.Repo.GetURL()
	case *github.PullRequestEvent:
		source = event.PullRequest.GetURL()
	case *github.PullRequestReviewEvent:
		source = event.PullRequest.GetURL()
	case *github.PullRequestReviewCommentEvent:
		source = event.Comment.GetURL()
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
	case *github.ReleaseEvent:
		source = event.Release.GetURL()
	case *github.RepositoryEvent:
		source = event.Repo.GetURL()
	case *github.StatusEvent:
		source = event.Commit.GetURL()
	case *github.TeamEvent:
		source = event.Team.GetURL()
	case *github.TeamAddEvent:
		source = event.Team.GetURL()
	case *github.WatchEvent:
		source = event.Repo.GetURL()
	}

	if source == "" {
		return nil, errors.New("unsupported event type")
	}

	s, err := url.Parse(source)
	if err != nil {
		return nil, err
	}

	ce := &cloudevents.Event{
		ID:              github.DeliveryID(req),
		Type:            fmt.Sprintf("com.github.%s", webHookType),
		Source:          *s,
		DataContentType: "application/json",
	}

	return ce, nil
}
