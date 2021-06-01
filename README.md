# Mastro

👷 Data and Feature Catalogue in Go 

## Data What?
![ML Process](img/ml_dev_process.png)

## Goals
Managing Metadata to achieve:
- Lineage - understanding data dependencies
- Quality - enrich data with dependability information
- Democratization - foster self-service culture

Modeling of the data transformation layer to foster:
- Reuse - of domain specific knowledge and pre-computed features in different use cases
- Enforcing - of established industry or domain-specific transformations and practices
- Interaction - between different roles through teams and projects

## Disclaimer

Mastro is still on development and largely untested. Please fork the repo and extend it at wish.

## TL-DR

Terminology:
* [Connector](doc/CONNECTORS.md) - component handling the connection to volumes and data bases
* [FeatureStore](doc/FEATURESTORE.md) - service to manage features (i.e., featureSets and featureStates);
* [Catalogue](doc/CATALOGUE.md) - service to manage data assets (i.e., static data definitions and their relationships);
* [Crawler](doc/CRAWLERS.md) - any agent able to list and walk a file system, filter and parse asset definitions (i.e. manifest files) and push them to the catalogue;

Help:
* [PlantUML Diagram of the repo](https://www.dumels.com/diagram/2e5f820a-1822-4852-8259-4811deefa789)
* [Configuration](doc/CONFIGURATION.md)
* [Deploy to K8s](doc/K8S-DEPLOY.md)

Status:

[![Go Build](https://github.com/data-mill-cloud/mastro/actions/workflows/go-build.yml/badge.svg)](https://github.com/data-mill-cloud/mastro/actions/workflows/go-build.yml)
[![Docker Image CI](https://github.com/data-mill-cloud/mastro/actions/workflows/docker-image.yml/badge.svg?branch=main)](https://github.com/data-mill-cloud/mastro/actions/workflows/docker-image.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
