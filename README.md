# cloudevents-webhook-gateway

cloudevents-webhook-gateway is a HTTP gateway. This receives the webhook requests from the service, converts it to CloudEvents format, and forwards the request to the backend specified in the configuration.

## Install

Download the latest binary from the [Releases](https://github.com/summerwind/cloudevents-webhook-gateway/releases) page.

Docker images are also available. Running cloudevents-webhook-gateway with Docker is as follows.

```
$ docker run -it -v $PWD/config.yml:/config.yml -p 24381:24381 summerwind/cloudevents-webhook-gateway:latest
```

## Usage

cloudevents-webhook-gateway can be started from the command line as follows.

```
$ cloudevents-webhook-gateway -c config.yml
```

To start cloudevents-webhook-gateway, specify the configuration file using the `-c` option. The configuration format is in YAML. Please see `example/config.yml` for the full configuration file format.

## Supported webhook

cloudevents-webhook-gateway currently supports the following webhooks.

- Github
- Docker Hub
- Alertmanager
- Anchore Engine
- Clair
- Slack
