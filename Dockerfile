FROM golang
MAINTAINER travis.simon@nicta.com.au

# Copy across our src files
ADD . /go/src/github.com/travissimon/microservices/uniq

# Build server
RUN go install github.com/travissimon/microservices/uniq

WORKDIR /go/bin

CMD uniq

# Listen on 8080
EXPOSE 8080