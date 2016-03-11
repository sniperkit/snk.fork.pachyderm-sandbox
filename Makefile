CLUSTER_NAME = pachyderm
PROJECT_NAME = pachyderm-sandbox
REGION = us-central1-a

# For docker publishing:
REPO=pachyderm/sandbox
TAG=latest

run:
	GIN_MODE=debug ./sandbox

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
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GO15VENDOREXPERIMENT=0 go get -d -v -t -u -f ./src/... ./app/...

vendor-without-update:
	go get -v github.com/kardianos/govendor
	rm -rf vendor
	govendor init
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 govendor add +external
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 govendor update +vendor
	$(foreach vendor_dir, $(VENDOR_IGNORE_DIRS), rm -rf vendor/$(vendor_dir) || exit; git checkout vendor/$(vendor_dir) || exit;)

vendor-for-google-app-engine: vendor
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 govendor add +local

vendor: vendor-update vendor-without-update

build:
	GO15VENDOREXPERIMENT=1 go build .

docker-build:
	docker build -t sandbox .

docker-debug:
	docker run --publish 9080:9080 sandbox

docker-push:
	docker login -e "$$DOCKER_EMAIL" -u "$$DOCKER_USERNAME" -p "$$DOCKER_PASSWORD"
	docker build -f Dockerfile -t $REPO:$$COMMIT .
	docker tag $REPO:$$COMMIT $REPO:$TAG
	docker tag $REPO:$$COMMIT $REPO:travis-$$TRAVIS_BUILD_NUMBER
	docker push $REPO

kube-generate-credentials:
	gcloud container clusters get-credentials pachyderm
	cp ~/.kube/config kube-config
	travis encrypt-file kube-config --add

kube-deploy:
	kubectl --kubeconfig="./kube-config" rolling-update sandbox --image=pachyderm/sandbox

deploy: docker-build docker-push kube-deploy

ci-setup:
	gcloud container clusters get-credentials pachyderm
	# generates update to ~/.kube/config