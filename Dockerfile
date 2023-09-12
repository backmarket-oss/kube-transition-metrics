FROM scratch

WORKDIR /

ARG BINARY
COPY $BINARY kube-transition-metrics /

EXPOSE 8080
USER 1000:1000
ENTRYPOINT ["/kube-transition-metrics"]
