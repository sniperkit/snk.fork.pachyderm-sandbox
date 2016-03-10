CLUSTER_NAME = pachyderm
PROJECT_NAME = pachyderm-sandbox
REGION = us-central1-a

run:
	go run sandbox.go

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
