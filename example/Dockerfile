FROM golang:1.17 AS build

RUN go get github.com/golang/dep/cmd/dep

COPY .  /go/src/github.com/[username]/[team-name]

WORKDIR /go/src/github.com/[username]/[team-name]

RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/team-name

FROM scratch
COPY --from=build /bin/team-name /bin/team-name
ENTRYPOINT ["/bin/team-name"]
CMD []
