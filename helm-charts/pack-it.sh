#!/usr/bin/env bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

echo "Packaging helm chart for mastro"
helm package mastro
echo "moving package at ../docs/helm-charts/"
mv mastro*.tgz ../docs/helm-charts/

cd ../docs/
helm repo index helm-charts --url https://data-mill-cloud.github.com/mastro/helm-charts
cd $SCRIPT_DIR