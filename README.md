# pvm-controller

This is a very simple Preemptible Controller, deployed to a GKE cluster it will monitor the Preemptible nodes used in the cluster and drain and delete nodes gracefully that are created close to each other, to prevent large disruptions when GKE automatically kills the nodes after 24h.

The controller is under development and we plan to release it as a proper helm chart once we reach version `1.0`.

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