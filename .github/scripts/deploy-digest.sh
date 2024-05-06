#!/bin/bash

IMAGE=$1
HELM_FILE=$2
HELM_PATH=infrastructure/argocd/projects/apps/$HELM

if [ -z "$HELM_FILE" ]; then
  echo "No deployment is required since there is no manifest"
  exit 0
fi

DIGEST=$(echo sha256:"${IMAGE##*@sha256:}")
echo "Digest is $DIGEST"

yq 'setpath(["spec", "source", "helm", "valuesObject", "image", "digest"]; "'$DIGEST'")' $HELM_PATH > $HELM_PATH
echo 'Patched ArgoCD manifests'
