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
RUN go build -o ./bin/mailroom mailroom

# Move bin to another image.
FROM alpine
WORKDIR /usr/src/app
COPY --from=builder /go/bin ./bin

# Add Add SSL ca certificates.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Add other files such as sql templates or jobs.
COPY ./sql ./sql
COPY ./jobs ./jobs

RUN chmod -R 755 /usr/src/app/jobs
RUN echo "Done."
