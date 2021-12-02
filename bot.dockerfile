FROM golang:1.17.3-alpine3.14 AS base
WORKDIR $GOPATH/src/gitlab.com/dreamteam-hack/hack041221/telegram-bot
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
RUN apk add --update --no-cache git ca-certificates gcc libc-dev make
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build-bot

FROM alpine:3.13.2
WORKDIR /app/
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group
COPY --from=base /go/src/gitlab.com/dreamteam-hack/hack041221/telegram-bot/bin/app /app/app
USER appuser:appuser
CMD ["/app/app"]
