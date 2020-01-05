---
title: "Writing Your Own Operator"
date: 2019-12-30T11:42:00Z
weight: 9
draft: true
---

There are more articles and tools to help with developing Operators, including using [Helm charts](https://docs.okd.io/latest/operators/osdk-helm.html). Based on nothing other than my own personal prejudices, I opted to try writing an Operator using Ansible and Golang. Both approaches use the [Operator SDK](https://github.com/operator-framework/operator-sdk), which is part of the [Operator Framework](https://coreos.com/blog/introducing-operator-framework). It can be installed as described at [Install Operator SDK](https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md).

Operators manage Custom Resources. On their own, Custom Resources simply let you store and retrieve structured data of a resource. The cluster still needs a custom controller to monitor the state of the resource and reconcile it to match the specified configuration. Together, a Custom Resource and its controller provides a declarative API.

In lieu of having sufficient imagination to build something useful, I decided to knock up an Operator to report when Experimentals, a contrived Custom Resource representing an entirely hypothetical service for storing experiment results in a blockchain and analysing them using AI, were created or changed in a k8s cluster. Let me be clear, it's only ever going to be the name of a Resource for demonstration purposes. Nothing clever is happening here. A k8s-based blockchain demonstration is available from IBM's [blockchain-network-on-kubernetes](https://github.com/IBM/blockchain-network-on-kubernetes).

Google has provided a set of [best practices](https://cloud.google.com/blog/products/containers-kubernetes/best-practices-for-building-kubernetes-operators-and-stateful-apps) for building Operators. It includes enabling logging and monitoring of the Operator. As my Operator is so trivial, I didn't attempt to follow these practices.

For readers who don't want to wade through the guff below, the executive summary is that it was relatively simple to use the Ansible Operator, and (as you'd expect) there's a bit more work getting the go-based Operator to work. I don't think I learned much about how Operators work, but the exercise was a reasonable refresher on k8s terminology and showed me I would go about starting to write my own real-world Operator, and I did get to fix some bugs (and introduce some more). Unless I was working on a legacy/non-standard application or for a vendor, I can't imagine needing to write my own Operator though.

# Ansible

The [Ansible Operator](https://www.ansible.com/blog/ansible-operator) provides features to manage OpenShift-based clusters, but can also be used for non-OpenShift clusters. It is built upon the 
[Operator SDK](https://github.com/operator-framework/operator-sdk), which is part of the [Operator Framework](https://coreos.com/blog/introducing-operator-framework). Its useful if you want to integrate with any other Ansible resources or playbooks, or if you don't have need or want to use more generic programming features.

A good starting place to learn about the Ansible Operator is the [Ansible Operator course](https://learn.openshift.com/ansibleop/ansible-operator-overview/?extIdCarryOver=true&sc_cid=701f2000001OH7YAAW). However, it's based on OpenShift, so the [Ansible Operator User Guide](https://github.com/operator-framework/operator-sdk/blob/master/doc/ansible/user-guide.md) may be more helpful for non-OpenShift users.

The Ansible Operator can trigger either a playbook or a role when the state of a Custom Resource changes. I chose to trigger a role, based on the conceit of an Experimental being a distinct kind of entity that would be managed in a common way, rather than just a set of related tasks.

An example of a Custom Resource used to specify the Experimental is as follows:

    apiVersion: "experiments.vapourware.com/v1alpha1"
    kind: "Experimental"
    metadata:
        name: "experimental"
    annotations:
        ansible.operator-sdk/reconcile-period: "30s"
    spec:
        is_public: True

The `reconcile-period` annotation above indicates that the k8s cluster should ensure the Experimental resources are made to match the specified configuration every 30 seconds.

The `spec` block also specifying parameters than can be used to customise the Ansible logic.

As is often the case with Ansible code, it is important the Operator code is idempotent.

## Generate the Operator scaffold code

The following command generates the Operator scaffold code:

    operator-sdk new experimental-operator --type=ansible --api-version=experiments.vapourware.com/v1alpha1 --kind=Experimental

It produces the structure shown below:

    experimental-operator
    ├── build
    │   ├── Dockerfile
    │   └── test-framework
    │       ├── Dockerfile
    │       └── ansible-test.sh
    ├── deploy
    │   ├── crds
    │   │   ├── experiments.vapourware.com_experimentals_crd.yaml
    │   │   └── experiments.vapourware.com_v1alpha1_experimental_cr.yaml
    │   ├── operator.yaml
    │   ├── role.yaml
    │   ├── role_binding.yaml
    │   └── service_account.yaml
    ├── molecule
    │   ├── default
    │   │   ├── asserts.yml
    │   │   ├── molecule.yml
    │   │   ├── playbook.yml
    │   │   └── prepare.yml
    │   ├── test-cluster
    │   │   ├── molecule.yml
    │   │   └── playbook.yml
    │   └── test-local
    │       ├── molecule.yml
    │       ├── playbook.yml
    │       └── prepare.yml
    ├── roles
    │   └── experimental
    │       ├── README.md
    │       ├── defaults
    │       │   └── main.yml
    │       ├── files
    │       ├── handlers
    │       │   └── main.yml
    │       ├── meta
    │       │   └── main.yml
    │       ├── tasks
    │       │   └── main.yml
    │       ├── templates
    │       └── vars
    │           └── main.yml
    └── watches.yaml

The mildly interesting/non-obvious directories/files are as follows:

* `build` - contains scripts that the operator-sdk uses for build and initialization
* `deploy` - contains a generic set of k8s manifests for deploying this Operator on a k8s cluster
* `watches.yaml` - maps the Custom Resource (identified by Group, Version, and Kind) to the Ansible code to invoke

*NB. I moved the `deploy/crd/experiments.vapourware.com_experimentals_crd.yaml` file to `deploy/crd/experiments.vapourware.com_v1alpha1_experimental_crd.yaml` to better fit with the examples I was following.*

## Define the Operator logic

All my Operator logic does is to deploy an Experimental resource.

BTW the watches.yml can be setup to invoke a role taken from Ansible Galaxy. To install a role from Galaxy, run the command shown below:

    ansible-galaxy install ROLE_NAME -p ./roles

My role is built upon the scaffold code generated previously. It uses the [k8s](https://docs.ansible.com/ansible/latest/modules/k8s_module.html) module to manage the Experimental k8s objects.

## Register the Experimental CRD

I created my cluster using [EKS](https://aws.amazon.com/eks/?nc2=h_ql_prod_cp_eks) and then deployed the Custom Resource Definition (CRD) on this cluster using the command shown below:

    kubectl create -f deploy/crds/experiments.vapourware.com_v1alpha1_experimental_crd.yaml

The CRD defines the type of the Custom Resource object.

## Running the Operator

The Operator can be run either as a Pod inside the k8s cluster, or as a go program outside the cluster using operator-sdk. I've used the first approach (which is recommended for production) with the following command:

    operator-sdk build experimental-operator:v0.0.1

I'll use the second approach in the Golang section.

The scaffold code in `deploy/operator.yml` contains a placeholder for the name of the previously built image. It can be updated as follows:

    sed -i 's|{{ REPLACE_IMAGE }}|YOUR_AWS_ACCOUNT_ID.dkr.ecr.eu-west-2.amazonaws.com/experimental-operator:v0.0.1|g' deploy/operator.yaml

I created an ECR repository, then tagged and deployed the image to ECR using the usual ECR/Docker commands. However, for k8s to pull the image from ECR, the following command is needed to sett the `imagePullPolicy` to `Always`:

    sed -i 's|{{ pull_policy\|default('\''Always'\'') }}|Always|g' deploy/operator.yaml

*NB. On MacOS, the two commands are as follows:*

    sed -i "" 's|{{ REPLACE_IMAGE }}|YOUR_AWS_ACCOUNT_ID.dkr.ecr.eu-west-2.amazonaws.com/experimental-operator:v0.0.1|g' deploy/operator.yaml
    sed -i "" 's|{{ pull_policy\|default('\''Always'\'') }}|Always|g' deploy/operator.yaml

If you don't want to publish your image to a registry, the `imagePullPolicy` should be set to `Never`. The Linux and MacOs commands below will make this change:

    sed -i "s|{{ pull_policy\|default('Always') }}|Never|g" deploy/operator.yaml
    sed -i "" 's|{{ pull_policy\|default('\''Always'\'') }}|Never|g' deploy/operator.yaml

To actually deploy the Operator, use the commands below:

    kubectl create -f deploy/service_account.yaml
    kubectl create -f deploy/role.yaml
    kubectl create -f deploy/role_binding.yaml
    kubectl create -f deploy/operator.yaml

*NB. role.yaml and role_binding.yaml describe cluster-wide resources. Creating these requires elevated permissions.*

Verifying that the Operator is running is done by checking the output of the following command:

    kubectl get deployment

## Create an Experimental CR

To create an Experimental Custom Resource instance using the Operator, set `deploy/crds/experiments.vapourware.com_v1alpha1_experimental_cr.yaml` as shown below:

    apiVersion: "experiments.vapourware.com/v1alpha1"
    kind: "Experimental"
    metadata:
        name: "experimental"
    spec:
        is_public: True
        number_of_replicas: 3

Then run the following command:

    kubectl apply -f deploy/crds/experiments.vapourware.com_v1alpha1_experimental_cr.yaml

Verify that the Operator and Experimental deployment is running (with the specified number of replicas) by checking the output of the following command:

    kubectl get deployment
    kubectl get pods

## Cleanup

To remove the Custom Resource instance, the CRD, and the Operator, run the following commands:

    kubectl delete -f deploy/crds/experiments.vapourware.com_v1alpha1_experimental_cr.yaml
    kubectl delete -f deploy/operator.yaml
    kubectl delete -f deploy/role_binding.yaml
    kubectl delete -f deploy/role.yaml
    kubectl delete -f deploy/service_account.yaml
    kubectl delete -f deploy/crds/experiments.vapourware.com_v1alpha1_experimental_crd.yaml

*NB. If you no longer need the k8s cluster you used for these deployments, tear that down too!*

## Debugging

To view the logs from the Ansible run, use the command below:

    kubectl logs deployment/experimental-operator -c ansible

To view more detailed logs from the Ansible run (including the Ansible Operator internals and the interface with k8s), use the command below:

    kubectl logs deployment/experimental-operator -c operator

# Golang

The article [Writing your first K8s operator](https://medium.com/faun/writing-your-first-kubernetes-operator-8f3df4453234) describes the steps involved in writing an Operator using go.

## Generate the Operator scaffold code

If you've set your working directory to be $GOPATH/src, the following command generates the Operator scaffold code:

    operator-sdk new experimental-operator

I have not set up a GOPATH, so had to use the following:

    operator-sdk new experimental-operator --repo experimental-operator

It produces the structure shown below:

    experimental-operator
    ├── build
    │   ├── Dockerfile
    │   └── bin
    │       ├── entrypoint
    │       └── user_setup
    ├── cmd
    │   └── manager
    │       └── main.go
    ├── deploy
    │   ├── operator.yaml
    │   ├── role.yaml
    │   ├── role_binding.yaml
    │   └── service_account.yaml
    ├── go.mod
    ├── go.sum
    ├── pkg
    │   ├── apis
    │   │   └── apis.go
    │   └── controller
    │       └── controller.go
    ├── tools.go
    └── version
        └── version.go

The mildly interesting/non-obvious directories/files are as follows:

* `build` - contains Docker code to build the Docker image
* `cmd/manager/main.go` - go code to run the Operator
* `deploy` - contains a generic set of k8s manifests for deploying this Operator on a k8s cluster

## Define the Operator logic

To add a CRD for the Experimental Operator, run the command below from the `experimental-operator` directory:

    operator-sdk add api --api-version=experiments.vapourware.com/v1alpha1 --kind=Experimental

The above command should have produced the following files:

* A CRD and an example Custom Resource in the `deploy/crds` directory.
* `pkg/apis/app/v1alpha1/experimental_types.go` defines the structure of the ExperimentalSpec which is the expected state of the Experiment Custom Resource and which is specified by the user in the `deploy/crds/experiments.vapourware.com_v1alpha1_experimental_cr.yaml` file. It also defines the structure of the ExperimentalStatus which will be used to provide the observed state of the Experimental when the kubectl describe command is executed.

*NB. Again, I moved the `deploy/crd/experiments.vapourware.com_experimentals_crd.yaml` file to `deploy/crd/experiments.vapourware.com_v1alpha1_experimental_crd.yaml` to match what I had done in the Ansible section. This turned out to be pretty dumb.*

To add the associated Controller for the Experimental Operator, run the command below from the `experimental-operator` directory:

    operator-sdk add controller --api-version=experiments.vapourware.com/v1alpha1 --kind=Experimental

This command above should have produced the file `pkg/controller/experimental/experimental_controller.go` (and the init code at `pkg/controller/add_experimental.go`) which should be customised to implement the business logic required. This logic typically takes the form of responding the changes in primary resource (`Experimental`) and any secondary resources. In this example, the second resource would probably be a PersistentVolume or Pod - changes in which could hypothetically require custom steps to be performed on the Experimental resources, say backing up data or updating the code in the Pods to handle changes in the blockchain structure. Yes, I have no idea how blockchain applications work.

[HERE](HERE) shows all the changes I made to the scaffold code.

I'm not going to attempt to write any real business logic. As a novice Go developer, it looks like I could work out how to fill in the placeholders with something sensible if I had real requirements, but it feels like it would take much longer than using Ansible. As it stands, the scaffold code does deploy an Experimental resource, which is all the Ansible Overseer example achieved anyway.

Out of the box, the Controller code uses the K8s API to create a single Pod if none with a given `app` label already exists. I'll just tweak it to add IsPublic and NumberOfReplicas fields to the spec, and  and status structures.

    type ExperimentalSpec struct {
      IsPublic bool `json:"isPublic"`
      NumberOfReplicas int32 `json:"numberOfReplicas"`
    }

    type ExperimentalStatus struct {
      NumberOfReplicas int32 `json:"numberOfReplicas"`
      PersistentVolumesUsed []string `json:"persistentVolumesUsed"`
    }

After adding the fields, the command below must be run:

    operator-sdk generate k8s

Next, the code would be updated to configure the primary and secondary resources that the Controller will monitor in the namespace. For this demonstration, this is already implemented as, by default, the Controller operates on a Experimental resource and creates a Pod (which can be assumed to be a reasonable secondary resource, if I don't want to put in any actual effort to this thing).

Lastly, the logic of scaling up and down the Pods and updating the Custom Resource status with the IDs of Persistent Volumes used would be implemented. All of this happens in the `Reconcile` function of the controller. The function is invoked each time the Experimental Custom Resource is changed or a change happens in the Pods belonging to the Experimental resource. A summary of what is happening during reconciliation is as follows:

1. The Controller fetches the Experimental resource in the current namespace and compares the value of its NumberOfReplicas field with the actual number of Pods that match the `app` label to decide whether Pods need to be created or deleted.
2. The Experimental resource status is updated with the correct field values for NumberOfReplicas and PersistentVolumesUsed (though I used fake values rather than fetching the PersistentVolume information required).
3. If Pods need to be added or removed, the `Reconcile` function adds or removes one pod at a time, returns, and waits for the next invocation (since it will be called after a pod was created or deleted). Any new pods are marked as “owned” by the Experimental primary resource using the controllerutil.SetControllerReference() function. Having this ownership in place means that when the Experimental resource is deleted, all its “child” pods are deleted as well.

My `Reconcile` code is almost entirely whipped from the [aforementioned tutorial's GitHub](https://github.com/xcoulon/podset-operator/blob/master/pkg/controller/podset/podset_controller.go#L86) - just with added errors.

## Running the Operator

The Operator was then run as a go program outside the cluster using operator-sdk via the following command:

    operator-sdk build experimental-operator:v0.0.1

The scaffold code in `deploy/operator.yml` contains a placeholder for the name of the previously built image. It can be updated as follows:

    sed -i 's|REPLACE_IMAGE|YOUR_AWS_ACCOUNT_ID.dkr.ecr.eu-west-2.amazonaws.com/experimental-operator:v0.0.1|g' deploy/operator.yaml

*NB. On MacOS, the appropriate sed is as follows:*

    sed -i "" 's|REPLACE_IMAGE|YOUR_AWS_ACCOUNT_ID.dkr.ecr.eu-west-2.amazonaws.com/experimental-operator:v0.0.1|g' deploy/operator.yaml

I created an ECR repository, then tagged and deployed the image to ECR using the usual ECR/Docker commands.

I then created my cluster using [EKS](https://aws.amazon.com/eks/?nc2=h_ql_prod_cp_eks) and deployed the Operator with the commands below:

    kubectl create -f deploy/service_account.yaml
    kubectl create -f deploy/role.yaml
    kubectl create -f deploy/role_binding.yaml
    kubectl create -f deploy/crds/experiments.vapourware.com_v1alpha1_experimental_crd.yaml
    kubectl create -f deploy/operator.yaml

*NB. role.yaml and role_binding.yaml describe cluster-wide resources. Creating these requires elevated permissions.*

To verify that the Operator is running, check the output of the following command:

    kubectl get deployment

## Create an Experimental CR

To create an Experimental Custom Resource instance using the Operator, add `deploy/crds/experiments.vapourware.com_v1alpha1_experimental_cr.yaml` with the contents below:

    apiVersion: "experiments.vapourware.com/v1alpha1"
    kind: "Experimental"
    metadata:
        name: "experimental"
    spec:
        isPublic: True
        numberOfReplicas: 3

Then run the following command:

    kubectl apply -f deploy/crds/experiments.vapourware.com_v1alpha1_experimental_cr.yaml

Verify that the Operator and Experimental deployment is running (with the specified number of replicas) by checking the output of the following command:

    kubectl get deployment
    kubectl get pods

Scaling can be tested by editing the `numberOfReplicas` field in `deploy/crds/experiments.vapourware.com_v1alpha1_experimental_cr.yaml` and running the commands below:

    kubectl apply -f deploy/crds/experiments.vapourware.com_v1alpha1_experimental_cr.yaml
    kubectl get pods

To check the custom status output, use the following command:

    kubectl describe experimental/example-experimental

## Cleanup

To remove the Custom Resource instance, the CRD, and the Operator, run the following commands:

    kubectl delete -f deploy/crds/experiments.vapourware.com_v1alpha1_experimental_cr.yaml
    kubectl delete -f deploy/operator.yaml
    kubectl delete -f deploy/role_binding.yaml
    kubectl delete -f deploy/role.yaml
    kubectl delete -f deploy/service_account.yaml
    kubectl delete -f deploy/crds/experiments.vapourware.com_v1alpha1_experimental_crd.yaml

*NB. If you no longer need the k8s cluster you used for these deployments, tear that down too!*

# Real World Example

See the [Tekton CD Operator](https://github.com/tektoncd/operator) for a more real-world example with an Operator being used to manage `Addon` and `Config` CRDs that manage the lifecycle of the created [Tekton CI/CD pipelines](https://github.com/tektoncd/pipeline). The video at [How to Build Cloud-Native CI/CD Pipelines with Tekton on Kubernetes](https://devfest.gdgnantes.com/sessions/how_to_build_cloud_native_ci_cd_pipelines_with_tekton_on_kubernetes/) discusses how to use Tekton Pipelines and briefly explains how the Operator is used to manage them.
