export PROJECT=$$(gcloud config get-value project)
export default-repo=gcr.io/${PROJECT}
export zone=us-west1-a
export cluster=beam
export TAG=v0.0.1

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build containers
	skaffold build --default-repo=${default-repo}

config: ## Deploy environment configuration to cluster
	kubectl create configmap job-config \
		--from-literal=INPUT="gs://apache-beam-samples/shakespeare/*" \
		--from-literal=OUTPUT="gs://${PROJECT}-output/wordcount/out"

deploy: ## Deploy artifacts
	skaffold run --default-repo=${default-repo}

connect: ## Connect to kubernetes cluster
	gcloud container clusters get-credentials ${cluster} --zone ${zone} --project ${project}