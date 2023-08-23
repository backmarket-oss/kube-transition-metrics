# syntax=docker/dockerfile:1

ARG BUILDPLATFORM
FROM --platform=$BUILDPLATFORM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

ARG GOARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build -o /kube-transition-monitoring

ARG TARGETPLATFORM
FROM --platform=$TARGETPLATFORM debian:buster-slim

WORKDIR /

COPY --from=build-stage /kube-transition-monitoring /

EXPOSE 2112

USER nouser:nogroup

ENTRYPOINT ["/kube-transition-monitoring"]
