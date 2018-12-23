FROM scratch

COPY ./cloudevents-webhook-gateway /bin/cloudevents-webhook-gateway

ENTRYPOINT ["/bin/cloudevents-webhook-gateway"]
