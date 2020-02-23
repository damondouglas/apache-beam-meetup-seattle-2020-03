export project=$$(gcloud config get-value project)
export default-repo=gcr.io/${project}
export zone=us-west1-a
export cluster=beam
export TAG=v0.0.1

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build containers
	skaffold build --default-repo=${default-repo}

config: connect ## Deploy environment configuration to cluster
	kubectl delete configmap job-config; \
	kubectl create configmap job-config \
		--from-literal=INPUT="gs://${project}-input" \
		--from-literal=OUTPUT="gs://${project}-output"

stage-tokenizer: ## Stage files for simulator
	curl https://zenodo.org/record/3238718/files/redmed_lexicon.tsv?download=1 | \
	gsutil cp - gs://${project}-input/redmed_lexicon.tsv

tokenizer: connect ## Deploy artifacts
	skaffold -p tokenizer delete; \
	skaffold -p tokenizer run --default-repo=${default-repo}

connect: ## Connect to kubernetes cluster
	gcloud container clusters get-credentials ${cluster} --zone ${zone} --project ${project}