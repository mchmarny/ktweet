
# Assumes following env vars set
#  GCP_PROJECT - ID of your project
#  CLUSTER_ZONE - GCP Zone, ideally same as your Knative k8s cluster

.PHONY: test image mod sample
.DEFAULT_GOAL := mod

all: test image

# DEV
test:
	go test ./... -v

# BUILD

mod:
	go mod tidy
	go mod vendor

image: mod
	gcloud builds submit \
		--project ${GCP_PROJECT} \
		--tag gcr.io/${GCP_PROJECT}/kuser

sample: mod
	gcloud builds submit \
		--project knative-samples \
		--tag gcr.io/knative-samples/kuser

# DEPLOYMENT

source:
	kubectl apply -f source.yaml -n demo

clean:
	rm -R ./vendor
	kubectl delete -f source.yaml -n demo


