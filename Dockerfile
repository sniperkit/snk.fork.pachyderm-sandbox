FROM gcr.io/google_appengine/golang
ENV GO15VENDOREXPERIMENT 1
COPY . /go/src/app
RUN ls /go/src/app
#RUN go-wrapper install -tags appenginevm
#RUN make build
RUN GO15VENDOREXPERIMENT=1 go install app/app.go