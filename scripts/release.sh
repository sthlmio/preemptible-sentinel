#!/bin/bash

set -e

# Download and install Google Cloud SDK
curl -o /tmp/google-cloud-sdk.tar.gz https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-253.0.0-linux-x86_64.tar.gz
tar -xvf /tmp/google-cloud-sdk.tar.gz -C /tmp/
/tmp/google-cloud-sdk/install.sh -q
source /tmp/google-cloud-sdk/path.bash.inc

echo $GOOGLE_CLOUD_SERVICE_KEY | base64 --decode -i > ${HOME}/gcloud-service-key.json
gcloud auth activate-service-account --key-file ${HOME}/gcloud-service-key.json
export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/gcloud-service-key.json

# Download and install Helm
curl -o /tmp/helm.tar.gz https://storage.googleapis.com/kubernetes-helm/helm-v2.14.1-linux-amd64.tar.gz
tar -xvf /tmp/helm.tar.gz -C /tmp/
mv /tmp/linux-amd64/helm /usr/local/bin/helm
helm init \
	--client-only \
	--skip-refresh

make release VER=${TRAVIS_TAG:1:-1}