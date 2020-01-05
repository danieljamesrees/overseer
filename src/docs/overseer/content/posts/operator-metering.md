---
title: "Operator Metering"
date: 2020-01-03T11:30:00Z
weight: 9
draft: false
---
The [Operator Metering](https://operatorhub.io/operator/metering-upstream) Operator aims to provide more knowledge about the usage and costs of running and managing k8s-native applications (i.e. Operators, rather than just applications that run on k8s). It is available as open source. The Operator ties into the cluster’s CPU and memory reporting, and calculates Infrastructure-as-a-Service (IaaS) costs and customised metrics (e.g. licensing).

Real world applications of Operator Metering include:

* Cloud budgeting - Provides insight into how resources are being used.
* Cloud billing - Allows resource usage billing that directly reflects the costs associated with project's actual infrastructure usage, driving accountability for service costs.
* Telemetry/aggregation - Allows viewing of service usage and metrics across many namespaces or teams.

See [Introducing the Operator Framework Metering](https://coreos.com/blog/introducing-operator-framework-metering) for a more detailed overview information of the tool.

Operator Metering is a collection of a few components:

* A Pod which aggregates [Prometheus](https://www.prometheus.io) data and generates reports based on the collected usage information.
* The [Presto](https://github.com/starburstdata/presto) database and the [Hive](https://cwiki.apache.org/confluence/display/Hive/Home#Home-ApacheHive) data warehousing application, that are used by the Operator Metering Pod to perform queries on the collected usage data.

See [Operator Metering Architecture](https://github.com/operator-framework/operator-metering/blob/master/Documentation/metering-architecture.md) for further information.

It would be worth comparing the reports made by possible by Operator Metering with those available from hosted solutions such as EKS, [GKE](https://cloud.google.com/kubernetes-engine/), and [AKS](https://azure.microsoft.com/en-gb/services/kubernetes-service/). As Operator Metering originated from CoreOS/Red Hat, presumably it's designed to particularly help those running on-premises/hybrid k8s environments.

My initial interest in Operator Metering was to see whether it could be used to monitor "DevOps" metrics, such as deployment frequency and Pod lifetime, as the initial description I read of its features was a little ambiguous. However, it does look like the tool is more focused on costs. An Operator could be used to such monitor cluster-wide DevOps metrics, in the same way they do for cost metrics with Operator Metering, though further research suggests the tool [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics) looks like it might already address this use case.

kube-state-metrics is a simple service that listens to the k8s API server and generates metrics about the state of the objects (including Deployments and Pods). Of course, a high-level deployment metric could be set up by simply directly reporting deployment activities from the deployment tool/pipeline itself. It is run by deploying the kube-state-metrics image inside a k8s Pod which has a Service Account token with read-only access to the cluster. [Example manifests](https://github.com/kubernetes/kube-state-metrics/tree/master/examples/standard) are available. The Deployment manifest should be customised to indicate which metrics should be generated. Prometheus can then discover kube-state-metrics instances using a specific Prometheus scrape config that picks up both the cluster metrics endpoint and kube-state-metrics own telemetry metrics.

As an aside, the [k8s Metrics API](https://kubernetes.io/docs/tasks/debug-application-cluster/resource-metrics-pipeline/) provides access to resource usage metrics, including CPU activity and memory consumption.
