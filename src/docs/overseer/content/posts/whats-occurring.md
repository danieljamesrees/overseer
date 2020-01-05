---
title: "What's Occurring"
date: 2019-12-30T09:07:51Z
weight: 7
draft: false
---
Overseer is a [Learn in Public](https://www.swyx.io/writing/learn-in-public/) project to help teach me about [Kubernetes Operators](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) and metrics on Kubernetes. Or an excuse to play round with [Hugo](https://gohugo.io), if I'm honest. The content is focused on areas I'm most likely to encounter, namely running simple Spring-based stateless microservices and widely used DevOps tools on a hosted [EKS](https://aws.amazon.com/eks/) cluster.

Operators are software extensions to Kubernetes that make use of custom resources to manage applications and their components. The Operator pattern aims to automate the steps a human operator takes to manage a service or set of services. Kubernetes’ controllers let you extend a cluster’s behaviour without modifying the code of Kubernetes itself. Operators are clients of the Kubernetes API that act as controllers for a [Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/). These Custom Resources are extensions of the Kubernetes API, that aren't necessarily available in a default Kubernetes installation. See [https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#should-i-add-a-custom-resource-to-my-kubernetes-cluster](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#should-i-add-a-custom-resource-to-my-kubernetes-cluster) to understand whether a Custom Resource is appropriate for your installation.

And I'm now bored of typing Kubernetes, so (being latte-quaffing hipster and all that), I'll hereafter refer to it as K8s.

Operators can be used to analyse and monitor the use of a K8s cluster, as with the following examples:

* [Container Security](https://operatorhub.io/operator/container-security-operator) - for identifying vulnerabilities in containers hosted on the K8s cluster
* [Grafana](https://operatorhub.io/operator/grafana-operator) - for running the Grafana platform to analyse and visualise the operations of your applications/services
* [Litmus Chaos](https://operatorhub.io/operator/litmuschaos) - for chaos engineering
* [Operator Metering](https://operatorhub.io/operator/metering-upstream) - for usage reporting of applications/services hosted on the K8s cluster (see also [Introducing the Operator Framework Metering](https://coreos.com/blog/introducing-operator-framework-metering)
* [Prometheus](https://operatorhub.io/operator/prometheus) - for running the Prometheus open-source solution to store application/service metrics in a time-series database

Operator Metering is covered in more detail at [operator-metering](../operator-metering).

Operators can also be used to automate the management of stateful/non-K8s-native applications/services, e.g. taking backups and performing schema upgrades alongside code upgrades. The Operators below provide examples of such complex applications:

* [Jenkins](https://operatorhub.io/operator/jenkins-operator) - for managing, umm, Jenkins
* [Sonatype Nexus](https://operatorhub.io/operator/nexus-operator-hub) - for managing binary artefacts
* [HashiCorp Vault](https://operatorhub.io/operator/vault) - for managing secrets

Interested readers, and I'm sure there will be many, can go to [https://www.weave.works/blog/how-to-correctly-handle-db-schemas-during-kubernetes-rollouts](https://www.weave.works/blog/how-to-correctly-handle-db-schemas-during-kubernetes-rollouts) to understand how Flyway can be used by Java applications/services to perform database schema migrations. It's a little hand-wavey IMHO, but still the best article I found after a quick Google session and maybe I'm thinking this stuff is more complicated than it actually is. [https://www.slideshare.net/nfrankel/zerodowntime-deployment-with-kubernetes-springboot-flyway](https://www.slideshare.net/nfrankel/zerodowntime-deployment-with-kubernetes-springboot-flyway) provides some other information that may aid understanding.

For users, such as myself, who shouldn't be doing anything clever with K8s, there should be no need to write my own Operator. As the above links suggest, access to community operators is provided from [OperatorHub.io](https://operatorhub.io). However to aid my understanding, I've taken a stab at it [writing-your-own-operator](writing-your-own-operator). The [Service Catalog](https://kubernetes.io/docs/concepts/extend-kubernetes/service-catalog/) provides one way of enabling applications running in K8s clusters to use external managed software offerings, such as a datastore service offered by a cloud provider. The Catalog can be installed using [Helm](https://kubernetes.io/docs/tasks/service-catalog/install-service-catalog-using-helm/).

For AWS users, the following Operators may be relevant:

* [AWS S3 Operator](https://operatorhub.io/operator/awss3-operator-registry)
* [AWS Service Operator](https://operatorhub.io/operator/aws-service)
* [API Operator](https://operatorhub.io/operator/api-operator) - for exposing microservices as API, though AWS users may prefer to use [API Gateway](https://aws.amazon.com/api-gateway/) outside of the K8s cluster
