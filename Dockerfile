FROM golang:1.13.4

WORKDIR /go/src/git.maxset.io/server/knaxim

RUN go get -v github.com/google/go-tika/tika
RUN go get -v github.com/gorilla/mux
RUN go get -v github.com/badoux/checkmail
RUN go get -v github.com/jdkato/prose/tokenize
RUN go get -v go.mongodb.org/mongo-driver/mongo
RUN go get -v github.com/gorilla/handlers

COPY ./config/containerconfig.json conf.json
COPY dockerrun.sh run.sh
RUN chmod 0755 ./run.sh
COPY resource ./resource

COPY passentropy ./passentropy
COPY srverror ./srverror
COPY database ./database
COPY *.go ./

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["./run.sh"]
