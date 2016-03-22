FROM golang:1.6

ADD . /go/src/github.com/pachyderm/sandbox
RUN go install -v github.com/pachyderm/pachyderm/sandbox

ENTRYPOINT /go/bin/sandbox
EXPOSE 80