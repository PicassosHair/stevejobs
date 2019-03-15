FROM golang:alpine

WORKDIR $GOPATH

COPY . .

RUN go build -o ./bin/parser parser
RUN go build -o ./bin/postdiff postdiff
