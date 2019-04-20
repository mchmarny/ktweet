
# Assumes following env vars set
#  GCP_PROJECT - ID of your project
#  CLUSTER_ZONE - GCP Zone, ideally same as your Knative k8s cluster

.PHONY: test image mod sample
.DEFAULT_GOAL := mod

all: test image

# DEV
test:
	go test ./... -v

# CONFIG

secret:
	kubectl create secret generic ktweet-secrets -n demo \
		--from-literal=T_CONSUMER_KEY=${T_CONSUMER_KEY} \
		--from-literal=T_CONSUMER_SECRET=${T_CONSUMER_SECRET} \
		--from-literal=T_ACCESS_TOKEN=${T_ACCESS_TOKEN} \
		--from-literal=T_ACCESS_SECRET=${T_ACCESS_SECRET}

# BUILD

mod:
	go mod tidy
	go mod vendor

image: mod
	gcloud builds submit \
		--project ${GCP_PROJECT} \
		--tag gcr.io/${GCP_PROJECT}/ktweet

sample: mod
	gcloud builds submit \
		--project knative-samples \
		--tag gcr.io/knative-samples/ktweet:0.1.1

# DEPLOYMENT

source:
	kubectl apply -f source.yaml -n demo

clean:
	rm -R ./vendor
	kubectl delete -f source.yaml -n demo


