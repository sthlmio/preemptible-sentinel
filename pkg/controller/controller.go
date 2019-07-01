/*
MIT License

Copyright (c) 2019 sthlm.io

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package controller

import (
	"github.com/sirupsen/logrus"
	"github.com/sthlmio/pvm-controller/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type PreemptibleController struct {
	client kubernetes.Interface
}

func Start() {
	stopCh := make(chan struct{})
	defer close(stopCh)

	c := NewPreemptibleController()
	go c.Run(stopCh)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

func NewPreemptibleController() *PreemptibleController {
	var client kubernetes.Interface
	_, err := rest.InClusterConfig()
	if err != nil {
		client = utils.GetClientOutOfCluster()
	} else {
		client = utils.GetClient()
	}

	pc := &PreemptibleController{
		client: client,
	}

	return pc
}

func (pc *PreemptibleController) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()

	logrus.Info("Starting Preemptible Controller")

	// Check things every 10 second.
	go wait.Until(pc.ListNodes, 10*time.Second, stopCh)
	<-stopCh
	logrus.Info("Shutting down Preemptible Controller")
}

func (pc *PreemptibleController) ListNodes() {
	options := metav1.ListOptions{
		LabelSelector: "cloud.google.com/gke-preemptible=true",
	}

	nodes, err := pc.client.CoreV1().Nodes().List(options)

	if err != nil {
		logrus.Fatalf("Error listing nodes (skipping rearrange): %v", err)
		return
	}

	for _, n := range nodes.Items {
		logrus.WithFields(logrus.Fields{
			"node":              n.Name,
			"creationTimestamp": n.CreationTimestamp,
		}).Infof("listed")
	}
}