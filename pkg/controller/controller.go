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
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sthlmio/preemptible-sentinel/pkg/config"
	"github.com/sthlmio/preemptible-sentinel/pkg/utils"
	v1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"
)

type PreemptibleController struct {
	Client kubernetes.Interface
	Config config.Config
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
	var c kubernetes.Interface
	_, err := rest.InClusterConfig()
	if err != nil {
		c = utils.GetClientOutOfCluster()
	} else {
		c = utils.GetClient()
	}

	pc := &PreemptibleController{
		Client: c,
		Config: config.Get(),
	}

	return pc
}

func (pc *PreemptibleController) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()

	logrus.Info("Starting Preemptible Controller")

	// Check things every 10 minute
	go wait.Until(pc.Process, pc.Config.DurationInMinutes*time.Minute, stopCh)
	<-stopCh
	logrus.Info("Shutting down Preemptible Controller")
}

func (pc *PreemptibleController) Process() {
	nodes, err := pc.ListNodes()

	if err != nil {
		logrus.Errorf("Error listing nodes (skipping rearrange): %v", err)
		return
	}

	if len(nodes.Items) <= 0 {
		logrus.Infof("No preemptible nodes found in cluster")
		return
	}

	// Sort nodes by creation timestamp (ASC sort)
	sort.SliceStable(nodes.Items, func(i, j int) bool { return nodes.Items[i].CreationTimestamp.UTC().Before(nodes.Items[j].CreationTimestamp.UTC()) })

	lengthOfNodeSlice := len(nodes.Items)
	for i, node := range nodes.Items {
		if !utils.IsNodeReady(node.Status) {
			continue
		}

		nextIndex := 1 + i
		if nextIndex < lengthOfNodeSlice {
			nextNode := nodes.Items[nextIndex]

			if !utils.IsNodeReady(nextNode.Status) {
				continue
			}

			if nextNode.CreationTimestamp.UTC().Sub(node.CreationTimestamp.UTC()).Minutes() < pc.Config.DeleteDiffMinutes && time.Now().UTC().Sub(node.CreationTimestamp.UTC()).Minutes() > 60 {
				logrus.WithFields(logrus.Fields{
					"node": node.Name,
				}).Infof("processing node termination")

				patchBytes := []byte(fmt.Sprint(`{"spec":{"unschedulable":true}}`))
				if _, err := pc.Client.CoreV1().Nodes().Patch(node.Name, types.StrategicMergePatchType, patchBytes); err != nil {
					logrus.Errorf("failed to patch node: %v", err)
					continue
				}

				pods, err := pc.ListPods(node.Name)

				if err != nil {
					logrus.Errorf("error listing pods: %v", err)
					continue
				}

				kubeSystemPods, err := pc.ListKubeSystemPods(node.Name)

				if err != nil {
					logrus.Errorf("error listing kube-system pods: %v", err)
					continue
				}

				pc.ProcessPods(filterPods(pods, "DaemonSet"))

				// We process kube-system pods last to allow for time to flush logs etc
				pc.ProcessPods(filterPods(kubeSystemPods, "DaemonSet"))

				logrus.Infof("evicted all pods")

				if err := pc.Client.CoreV1().Nodes().Delete(node.Name, &metav1.DeleteOptions{}); err != nil {
					logrus.WithFields(logrus.Fields{
						"node": node.Name,
					}).Errorf("failed to delete node: %v", err)

					continue
				}

				logrus.WithFields(logrus.Fields{
					"node": node.Name,
				}).Infof("successfully deleted node")

				break
			}

			logrus.WithFields(logrus.Fields{
				"node": node.Name,
			}).Infof("node does not match the delete criteria")
		}
	}
}

func (pc *PreemptibleController) ProcessPods(pods []v1.Pod) {
	nextPod:
	for _, p := range pods {
		logrus.WithFields(logrus.Fields{
			"pod":       p.Name,
			"namespace": p.Namespace,
		}).Infof("trying to delete pod")

		if err := pc.Client.CoreV1().Pods(p.Namespace).Delete(p.Name, &metav1.DeleteOptions{}); err != nil {
			logrus.WithFields(logrus.Fields{
				"pod":       p.Name,
				"namespace": p.Namespace,
			}).Errorf("failed to delete pod: %v", err)
		} else {
			for _, ownerReference := range p.ObjectMeta.OwnerReferences {
				// Don't wait for statefulsets, since they will be recreated immediately with the same name
				// and don't pass next check we make
				if ownerReference.Kind == "StatefulSet" {
					logrus.WithFields(logrus.Fields{
						"pod":       p.Name,
						"namespace": p.Namespace,
					}).Infof("pod was successfully deleted")
					continue nextPod
				}
			}

			if err := pc.CheckIfPodIsDeleted(p); err != nil {
				logrus.WithFields(logrus.Fields{
					"pod":       p.Name,
					"namespace": p.Namespace,
				}).Errorf("pod did not get deleted: %v", err)
			} else {
				logrus.WithFields(logrus.Fields{
					"pod":       p.Name,
					"namespace": p.Namespace,
				}).Infof("pod was successfully deleted")
			}
		}
	}
}

func (pc *PreemptibleController) ListNodes() (*v1.NodeList, error) {
	options := metav1.ListOptions{
		LabelSelector: "cloud.google.com/gke-preemptible=true",
	}

	return pc.Client.CoreV1().Nodes().List(options)
}

func (pc *PreemptibleController) ListKubeSystemPods(nodeName string) (*v1.PodList, error) {
	options := metav1.ListOptions{
		FieldSelector: fields.AndSelectors(
			fields.OneTermEqualSelector("spec.nodeName", nodeName),
			fields.OneTermEqualSelector("metadata.namespace", "kube-system"),
		).String(),
	}

	return pc.Client.CoreV1().Pods(metav1.NamespaceAll).List(options)
}

func (pc *PreemptibleController) ListPods(nodeName string) (*v1.PodList, error) {
	options := metav1.ListOptions{
		FieldSelector: fields.AndSelectors(
			fields.OneTermEqualSelector("spec.nodeName", nodeName),
			fields.OneTermNotEqualSelector("metadata.namespace", "kube-system"),
		).String(),
	}

	return pc.Client.CoreV1().Pods(metav1.NamespaceAll).List(options)
}

func (pc *PreemptibleController) CheckIfPodIsDeleted(p v1.Pod) error {
	return wait.PollImmediate(time.Second, 60*time.Second, func() (bool, error) {
		_, err := pc.Client.CoreV1().Pods(p.Namespace).Get(p.Name, metav1.GetOptions{})
		if apierrs.IsNotFound(err) {
			return true, nil // done
		}

		if err != nil {
			return true, err // stop wait with error
		}

		return false, nil
	})
}

func filterPods(podList *v1.PodList, kind string) (output []v1.Pod) {
	for _, pod := range podList.Items {
		for _, ownerReference := range pod.ObjectMeta.OwnerReferences {
			if ownerReference.Kind != kind {
				output = append(output, pod)
			}
		}
	}

	return
}
