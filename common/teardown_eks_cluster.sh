#!/usr/bin/env bash

set -o nounset
set -o pipefail

eksctl delete cluster --name overseer

while eksctl get cluster --name=overseer > /dev/null; do
  echo
  echo Waiting for overseer cluster to die - $(date)
  echo
  sleep 3
done

echo
echo Ensure the ENI with the Security Groups eks-cluster-sg-overseer-* is deleted
echo Ensure the Security Groups eks-cluster-sg-overseer-* is deleted
echo Ensure the VPC eksctl-overseer-cluster is deleted
echo

aws cloudformation delete-stack --stack-name eksctl-overseer-cluster

while [[ $(aws cloudformation list-stacks --query "StackSummaries[?StackName=='eksctl-overseer-cluster']" --stack-status-filter CREATE_IN_PROGRESS CREATE_FAILED CREATE_COMPLETE ROLLBACK_IN_PROGRESS ROLLBACK_FAILED ROLLBACK_COMPLETE DELETE_IN_PROGRESS DELETE_FAILED UPDATE_IN_PROGRESS UPDATE_COMPLETE_CLEANUP_IN_PROGRESS UPDATE_COMPLETE UPDATE_ROLLBACK_IN_PROGRESS UPDATE_ROLLBACK_FAILED UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS UPDATE_ROLLBACK_COMPLETE REVIEW_IN_PROGRESS IMPORT_IN_PROGRESS IMPORT_COMPLETE IMPORT_ROLLBACK_IN_PROGRESS IMPORT_ROLLBACK_FAILED IMPORT_ROLLBACK_COMPLETE) != "[]" ]]; do
  echo
  echo Waiting for overseer CloudFormation stack to die - $(date)
  echo
  sleep 3
done
