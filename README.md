# pvm-controller

This is a very simple Preemptible Controller, deployed to a GKE cluster it will monitor the Preemptible nodes used in the cluster and drain and delete nodes gracefully that are created close to each other, to prevent large disruptions when GKE automatically kills the nodes after 24h.

The controller is under development and we plan to release it as a proper helm chart once we reach version `1.0`.