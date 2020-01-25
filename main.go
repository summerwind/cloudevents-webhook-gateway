package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	uuid "github.com/satori/go.uuid"
	yaml "gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
	"github.com/summerwind/cloudevents-webhook-gateway/config"
	"github.com/summerwind/cloudevents-webhook-gateway/proxy"
	"github.com/summerwind/cloudevents-webhook-gateway/webhook"
	"github.com/summerwind/cloudevents-webhook-gateway/webhook/alertmanager"
	"github.com/summerwind/cloudevents-webhook-gateway/webhook/anchoreengine"
	"github.com/summerwind/cloudevents-webhook-gateway/webhook/clair"
	"github.com/summerwind/cloudevents-webhook-gateway/webhook/dockerhub"
	"github.com/summerwind/cloudevents-webhook-gateway/webhook/github"
	"github.com/summerwind/cloudevents-webhook-gateway/webhook/slack"
)

var (
	VERSION = "0.0.1"
	COMMIT  = "HEAD"
)

// loadConfig loads the specified configuration file and returns
// config.
func loadConfig(configPath string) (*config.Config, error) {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	c := config.New()
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func newProxyHandler(backend *url.URL, parser webhook.Parser) (*httputil.ReverseProxy, error) {
	director := func(req *http.Request) {
		// Copy request body
		body := req.Body
		if body != nil && body != http.NoBody {
			var buf bytes.Buffer

			_, err := buf.ReadFrom(body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to read request body: %s\n", err)
				return
			}

			err = body.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
				return
			}

			body = ioutil.NopCloser(&buf)
			req.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
		}

		ce, err := parser.Parse(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %s\n", err)
			return
		}

		if ce.ID == "" {
			id := uuid.NewV4()
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to generate event ID: %s\n", err)
				return
			}
			ce.ID = id.String()
		}

		if ce.Time == nil {
			t := time.Now()
			ce.Time = &t
		}

		req.Body = body

		req.Host = backend.Host
		req.URL.Scheme = backend.Scheme
		req.URL.Host = backend.Host
		req.URL.Path = backend.Path

		req.Header.Set("ce-specversion", "1.0")
		req.Header.Set("ce-type", ce.Type)
		req.Header.Set("ce-source", ce.Source.String())
		req.Header.Set("ce-id", ce.ID)

		if ce.Subject != "" {
			req.Header.Set("ce-subject", ce.Subject)
		}
		if ce.Time != nil {
			req.Header.Set("ce-time", ce.Time.Format(time.RFC3339))
		}
		if ce.DataSchema.String() != "" {
			req.Header.Set("ce-dataschema", ce.DataSchema.String())
		}
		if ce.DataContentType != "" {
			req.Header.Set("Content-Type", ce.DataContentType)
		}

		log.Printf("remote_addr:%s event_id:%s event_type:%s source:%s", req.RemoteAddr, ce.ID, ce.Type, ce.Source.String())
	}

	transport := proxy.NewTransport()

	return &httputil.ReverseProxy{Director: director, Transport: transport}, nil
}

// run starts the HTTP server to process authentication.
func run(cmd *cobra.Command, args []string) error {
	v, err := cmd.Flags().GetBool("version")
	if err != nil {
		return err
	}

	if v {
		fmt.Printf("%s (%s)\n", VERSION, COMMIT)
		return nil
	}

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	c, err := loadConfig(configPath)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	if c.GitHub.Backend != "" {
		backend, err := url.Parse(c.GitHub.Backend)
		if err != nil {
			return err
		}
		parser := github.NewParser(c.GitHub.Secret)

		handler, err := newProxyHandler(backend, parser)
		if err != nil {
			return err
		}

		mux.Handle(c.GitHub.Path, handler)
	}

	if c.DockerHub.Backend != "" {
		backend, err := url.Parse(c.DockerHub.Backend)
		if err != nil {
			return err
		}
		parser := dockerhub.NewParser()

		handler, err := newProxyHandler(backend, parser)
		if err != nil {
			return err
		}

		mux.Handle(c.DockerHub.Path, handler)
	}

	if c.Alertmanager.Backend != "" {
		backend, err := url.Parse(c.Alertmanager.Backend)
		if err != nil {
			return err
		}
		parser := alertmanager.NewParser()

		handler, err := newProxyHandler(backend, parser)
		if err != nil {
			return err
		}

		mux.Handle(c.Alertmanager.Path, handler)
	}

	if c.AnchoreEngine.Backend != "" {
		backend, err := url.Parse(c.AnchoreEngine.Backend)
		if err != nil {
			return err
		}
		parser := anchoreengine.NewParser()

		handler, err := newProxyHandler(backend, parser)
		if err != nil {
			return err
		}

		mux.Handle(c.AnchoreEngine.Path, handler)
	}

	if c.Clair.Backend != "" {
		backend, err := url.Parse(c.Clair.Backend)
		if err != nil {
			return err
		}
		parser := clair.NewParser()

		handler, err := newProxyHandler(backend, parser)
		if err != nil {
			return err
		}

		mux.Handle(c.Clair.Path, handler)
	}

	if c.Slack.Backend != "" {
		backend, err := url.Parse(c.Slack.Backend)
		if err != nil {
			return err
		}
		parser := slack.NewParser()

		handler, err := newProxyHandler(backend, parser)
		if err != nil {
			return err
		}

		mux.Handle(c.Slack.Path, handler)
	}

	server := &http.Server{
		Addr:    c.Listen,
		Handler: mux,
	}

	go func() {
		if c.TLS.CertFile != "" {
			server.ListenAndServeTLS(c.TLS.CertFile, c.TLS.KeyFile)
		} else {
			server.ListenAndServe()
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var cmd = &cobra.Command{
		Use:   "cloudevents-webhook-gateway",
		Short: "A gateway that converts webhook requests to CloudEvents",
		RunE:  run,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().StringP("config", "c", "config.yml", "Path to the configuration file")
	cmd.Flags().BoolP("version", "v", false, "Display version information and exit")

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
