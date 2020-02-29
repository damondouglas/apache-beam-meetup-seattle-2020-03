# About
This is the repository for the 2020-03-02 Apache Beam Meetup.

# What problem does this solve?

There are situations where patient allergies are inputed into a chart as free-text.  Existing allergy warning detection systems require coding of medication allergies using standard coding schemes such as RxNorm or Snomed.  This repository illustrates the use of Apache Beam to match two datasets by a matching drug name key.

# Requirements

- [terraform](https://terraform.io)
- [gcloud](https://cloud.google.com/sdk/gcloud)
- [skaffold](https://skaffold.dev)
- [golang](https://golang.org/)
- [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/)

# Setup

## Google Cloud platform

```
$ cd artifacts
$ terraform init
$ terraform plan
$ terraform apply
```

## Container images

```
$ make build
```

# Usage

```
$ make pipeline
```

# Patient data simulation
Patient data was simulated using randomly generated Medical Record Numbers (MRN) and assigned drug names acquired from the [redmed published dataset](https://zenodo.org/record/3238718#.Xk9hOpNKjOQ) by Lavertu and Altman.

Lavertu and Altman. 2019 RedMed: Extending drug lexicons for social media applications. https://doi.org/10.1016/j.jbi.2019.103307 

# Future work

## Handle drug name misspellings

This implementation only matches exact drug names.  Future work should explore fuzzy matching as most free-text medical documentation contains spelling errors.

## Execute pipeline on non-direct runners

The current status of the Apache Beam Golang runtime allows for direct and Dataflow runners.  This implementation utilized a direct runner.  This is fine for small datasets.  Future work should explore running on other runners such as Apache Flink.