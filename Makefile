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
	kubectl delete configmap job-config; \
	kubectl create configmap job-config \
		--from-literal=INPUT="gs://${PROJECT}-input" \
		--from-literal=OUTPUT="gs://${PROJECT}-output"

stage-simlator: ## Stage files for simulator
	curl https://zenodo.org/record/3238718/files/redmed_lexicon.tsv?download=1 | \
	gsutil cp - gs://${PROJECT}-input/redmed_lexicon.tsv

deploy: ## Deploy artifacts
	skaffold delete; \
	skaffold run --default-repo=${default-repo}

connect: ## Connect to kubernetes cluster
	gcloud container clusters get-credentials ${cluster} --zone ${zone} --project ${project}