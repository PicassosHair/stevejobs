FROM golang:alpine as builder

ENV GOPATH /go
WORKDIR $GOPATH

COPY . .

# Install deps.
RUN go get github.com/mailgun/mailgun-go/v3

# Build go binary and output to ${GOPATH}/bin
RUN go build -o ./bin/parser parser
RUN go build -o ./bin/postdiff postdiff

# Move bin to another image.
FROM alpine
WORKDIR /usr/src/app
COPY --from=builder /go/bin ./bin
# Add other files such as sql templates or jobs.
COPY ./sql ./sql
COPY ./jobs ./jobs

RUN chmod -R 755 /usr/src/app/jobs
RUN echo "Done."
