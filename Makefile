export project=$$(gcloud config get-value project)
export default-repo=gcr.io/${project}
export zone=us-west1-a
export cluster=beam
export TAG=v0.0.1

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build containers
	skaffold build --default-repo=${default-repo}

stage-tokenizer: ## Stage files for simulator
	curl https://zenodo.org/record/3238718/files/redmed_lexicon.tsv?download=1 | \
	gsutil cp - gs://${project}-input/redmed_lexicon.tsv

tokenizer: connect ## Deploy artifacts
	kubectl delete configmap tokenizer-config; \
	kubectl create configmap tokenizer-config \
		--from-literal=INPUT="gs://${project}-input/redmed_lexicon.tsv" \
		--from-literal=OUTPUT="bigquery://${project}:beam.drugs" \
		--from-literal=PROJECT="${project}" \
		--from-literal=COLUMNS="1,3,4,5"; \
	skaffold -p tokenizer delete; \
	skaffold -p tokenizer run --default-repo=${default-repo}

pipeline: connect ## Deploy artifacts
	kubectl delete configmap pipeline-config; \
	kubectl create configmap pipeline-config \
		--from-literal=PATIENTS="bigquery://${project}:patients.uncoded_allergies_sample" \
		--from-literal=OUTPUT="bigquery://${project}:patients.coded_patients_sample" \
		--from-literal=PROJECT="${project}" \
		--from-literal=SNOMED="bigquery://${project}:rxnorm.snomed_sample"; \
	skaffold -p pipeline delete; \
	skaffold -p pipeline run --default-repo=${default-repo}

connect: ## Connect to kubernetes cluster
	gcloud container clusters get-credentials ${cluster} --zone ${zone} --project ${project}