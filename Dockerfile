FROM scratch

COPY ./cloudevents-gateway /bin/cloudevents-gateway

ENTRYPOINT ["/bin/cloudevents-gateway"]
