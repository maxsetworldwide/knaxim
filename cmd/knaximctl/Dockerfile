FROM golang:1.13.4 AS backend

WORKDIR /go/src/git.maxset.io/web/knaxim

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY pkg ./pkg
COPY internal ./internal
COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo -v ./cmd/knaximctl

FROM alpine:3.12 AS production

COPY --from=backend /go/bin/knaximctl /bin/knaximctl
COPY container.config.json /etc/knaxim/conf.json
COPY resource /resource
COPY dockerrun.sh /bin/run.sh
RUN chmod 0755 /bin/run.sh
CMD ["/bin/run.sh"]
