#!/usr/bin/env bash

set -o nounset
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cd ${__dir}/experimental-operator

kubectl delete -f deploy/crds/experiments.vapourware.com_v1alpha1_experimental_cr.yaml
kubectl delete -f deploy/service_account.yaml
kubectl delete -f deploy/role.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/operator.yaml
kubectl delete -f deploy/crds/experiments.vapourware.com_v1alpha1_experimental_crd.yaml
