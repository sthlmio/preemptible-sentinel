# preemptible-sentinel

[![Build Status](https://travis-ci.org/sthlmio/preemptible-sentinel.svg?branch=master)](https://travis-ci.org/sthlmio/preemptible-sentinel)

This is a very simple Preemptible Sentinel (Controller), deployed to a GKE cluster it will monitor the Preemptible nodes used in the cluster and drain and delete nodes gracefully that are created close to each other, to prevent large disruptions when GKE automatically kills the nodes after 24h.

The controller is under development and should not be considered stable until we reach version `1.0`.

#### We use this in conjuction with:
- [estafette-gke-node-pool-shifter](https://github.com/estafette/estafette-gke-node-pool-shifter)
    - Shifting node pool from the backup node pool to our preemptible node pool if the autoscaler creates these backup nodes for some reason
- [k8s-node-termination-handler](https://github.com/GoogleCloudPlatform/k8s-node-termination-handler)
    - Used when GKE terminates a preemptible node so the node gracefully evicts all pods before deletion

#### Usage example:
- **Static node pool** (2 or more nodes for special workloads)
    - Perfect for special workload like ingresses, to avoid nodes entering and leaving load balancers etc. Or other very important workloads that always needs to be running
- **Preemptible node pool** (5-7 nodes using autoscaling)
    - Running all HA workloads by default
- **Backup node pool** (0-7 nodes using autoscaling)
    - Backup nodes if preemptible nodes should be out of stock
    
#### Install
Add sthlmio chart repository before installing the chart. Also the chart is installed with `--devel` flag to allow semver versions like `0.1.0-alpha.0` until we reach stable `1.0.0`.
```bash
helm repo add sthlmio https://charts.sthlm.io

helm install \
    --name preemptible-sentinel \
    --namespace sthlmio \
    --devel \
    sthlmio/preemptible-sentinel
```

#### Development
The development of the chart can only be done against a GKE cluster with a node pool of regular vms and a node pool of preemptible vms.
We use `go1.12`, good commands to keep in mind:

```bash
go build
go test -v ./...
go mod tidy
go mod vendor
```

##### Fast local development against GKE cluster
Make sure that current context is against the cluster to test against
`kubectl config set-context $(kubectl config current-context)
```bash
make && ./controller
```

##### Build/push/deploy for local development
```bash
export PRIVATE_DOCKER_REPO=<your private docker repo>

docker build --no-cache -t $(PRIVATE_DOCKER_REPO):latest .
docker push $(PRIVATE_DOCKER_REPO):latest
helm install \
    --name preemptible-sentinel \
    --namespace sthlmio \
    --set-string repository=$(PRIVATE_DOCKER_REPO) \
    --set-string tag=latest \
    --set-string pullPolicy=Always \
    ./chart/preemptible-sentinel
```