FROM golang:1.6

ADD . /go/src/github.com/pachyderm/sandbox
RUN go install -v github.com/pachyderm/sandbox

# Change working directory for template file reads
WORKDIR "/go/src/github.com/pachyderm/sandbox"

ENTRYPOINT /go/bin/sandbox
EXPOSE 80