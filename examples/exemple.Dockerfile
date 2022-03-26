FROM golang:1.17 AS build

COPY .  /my-bot

WORKDIR /my-bot

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/my-bot main.go

FROM scratch
COPY --from=build /bin/team-name /bin/team-name
ENTRYPOINT ["/bin/team-name"]
CMD []
