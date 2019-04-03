FROM golang:alpine as builder

ENV GOPATH /go
WORKDIR $GOPATH

COPY . .

# Install deps.
RUN apk add --no-cache git
RUN go get github.com/mailgun/mailgun-go

# Build go binary and output to ${GOPATH}/bin
RUN go build -o ./bin/parser parser
RUN go build -o ./bin/postdiff postdiff
RUN go build -o ./bin/mail mail

# Move bin to another image.
FROM alpine
WORKDIR /usr/src/app

# Install deps.
RUN apk update
RUN apk add --no-cache --no-progress ca-certificates wget unzip mysql-client bash curl
RUN update-ca-certificates

# Mounting outside disk here.
RUN mkdir /data

COPY --from=builder /go/bin ./bin

# Install Slack CLI tool for msg.
RUN curl -o /usr/src/app/bin/slack https://raw.githubusercontent.com/rockymadden/slack-cli/master/src/slack
RUN chmod +x /usr/src/app/bin/slack

# Add other files such as sql templates or jobs.
COPY ./sql ./sql
COPY ./jobs ./jobs
# Mysql config file used to connect to the database.
COPY ./mysql.conf ./mysql.conf

RUN chmod -R 755 /usr/src/app/jobs
RUN echo "Done."
