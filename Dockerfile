FROM golang:1.12.3

ENV GO111MODULE=on
WORKDIR /go/src/github.com/nomasters/hashmap

# copy over files important to the project
COPY go.mod .
COPY *.go ./
COPY hashmap/ hashmap

# install the commandline tools
RUN go install github.com/nomasters/hashmap/hashmap
