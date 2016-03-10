CLUSTER_NAME = pachyderm
PROJECT_NAME = pachyderm-sandbox
REGION = us-central1-a

run:
	goapp serve app/

setup:
	gcloud config set compute/zone $(REGION)
	gcloud config set project $(PROJECT_NAME)

pachctl:
	go install github.com/pachyderm/pachyderm/src/cmd/pachctl

kubectl: setup
	gcloud config set container/cluster $(CLUSTER_NAME)
	gcloud container clusters get-credentials $(CLUSTER_NAME)
	gcloud components update kubectl

cluster: setup
	 gcloud container clusters create $(CLUSTER_NAME)

pachyderm: cluster kubectl pachctl
	pachctl manifest | kubectl create -f -

vendor-update:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GO15VENDOREXPERIMENT=0 go get -d -v -t -u -f ./src/... ./.

vendor-without-update:
	go get -v github.com/kardianos/govendor
	rm -rf vendor
	govendor init
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 govendor add +external
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 govendor update +vendor
	$(foreach vendor_dir, $(VENDOR_IGNORE_DIRS), rm -rf vendor/$(vendor_dir) || exit; git checkout vendor/$(vendor_dir) || exit;)

vendor: vendor-update vendor-without-update

build:
	GO15VENDOREXPERIMENT=1 go build -o sandbox app/app.go
