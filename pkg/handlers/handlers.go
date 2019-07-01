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

package handlers

import (
	"github.com/sirupsen/logrus"
	"github.com/sthlmio/pvm-controller/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Handler is implemented by any handler.
// The Handle method is used to process event
type Handler interface {
	Init() error
	ObjectCreated(node string)
	ObjectDeleted(node string)
	ObjectUpdated(node string)
}

// Map maps each event handler function to a name for easily lookup
var Map = map[string]interface{}{
	"default": &Default{},
}

// Default handler implements Handler interface
type Default struct {
}

// Init initializes handler configuration
// Do nothing for default handler
func (d *Default) Init() error {
	return nil
}

func (d *Default) ObjectCreated(node string) {

}

func (d *Default) ObjectDeleted(node string) {
	logrus.WithFields(logrus.Fields{
		"node": node,
	}).Infof("deleted")

	var kubeClient kubernetes.Interface
	_, err := rest.InClusterConfig()
	if err != nil {
		kubeClient = utils.GetClientOutOfCluster()
	} else {
		kubeClient = utils.GetClient()
	}

	nodes, err := kubeClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		logrus.Fatalf("Error listing nodes (skipping rearrange): %v", err)
	} else {
		for _, n := range nodes.Items {
			if n.Labels["cloud.google.com/gke-preemptible"] == "true" {
				logrus.WithFields(logrus.Fields{
					"node": n.Name,
					"creationTimestamp": n.CreationTimestamp,
				}).Infof("listed")
			}
		}
	}
}

func (d *Default) ObjectUpdated(node string) {

}
