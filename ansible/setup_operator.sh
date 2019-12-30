#!/usr/bin/env bash

if [[ -z "${AWS_ACCOUNT}" ]]; then
  echo Must specify AWS_ACCOUNT
  exit
fi

set -o errexit
set -o nounset
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

build_and_push_image_to_ecr() {
  operator-sdk build experimental-operator:v0.0.1

  if ! aws ecr describe-repositories --repository-names=experimental-operator; then
    aws ecr create-repository --repository-name experimental-operator
  fi

  docker tag experimental-operator:v0.0.1 ${AWS_ACCOUNT}.dkr.ecr.eu-west-2.amazonaws.com/experimental-operator:v0.0.1

  eval "$(aws ecr get-login --no-include-email)"

  docker push ${AWS_ACCOUNT}.dkr.ecr.eu-west-2.amazonaws.com/experimental-operator:v0.0.1
}

cd ${__dir}/experimental-operator

# Validation turned off due to a known issue with the field "x-kubernetes-preserve-unknown-fields".
kubectl create -f deploy/crds/experiments.vapourware.com_v1alpha1_experimental_crd.yaml --validate=false

build_and_push_image_to_ecr

kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/operator.yaml

kubectl get deployment

kubectl apply -f deploy/crds/experiments.vapourware.com_v1alpha1_experimental_cr.yaml

kubectl get deployment
kubectl get pods
