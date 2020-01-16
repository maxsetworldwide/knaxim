FROM golang:1.13.4

WORKDIR /go/src/git.maxset.io/web/knaxim

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go install -v ./...

COPY pkg ./pkg
COPY internal ./internal
COPY cmd ./cmd

RUN go install -v ./...

COPY container.config.json /etc/knaxim/conf.json
COPY resource ./resource
COPY dockerrun.sh run.sh
RUN chmod 0755 ./run.sh

CMD ["./run.sh"]
