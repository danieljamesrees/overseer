#!/usr/bin/env bash

# Only works on MacOS.

set -o errexit
set -o nounset
set -o pipefail

setup_eksctl() {
  brew tap weaveworks/tap

  if eksctl version; then
    brew upgrade eksctl && brew link --overwrite eksctl
  else
    eksctl version
  fi
}

setup_eksctl

# Fargate fails due to a permissions issue.
# 2 nodes are needed for a single experimental-operator pod.
# 3 nodes are needed for 2 additional single experimental-operator pods.
eksctl create cluster \
--name overseer \
--version 1.14 \
--nodegroup-name overseer-workers \
--node-type t2.micro \
--nodes 3 \
--nodes-min 3 \
--nodes-max 3 \
--managed
