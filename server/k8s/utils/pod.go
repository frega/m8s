package utils

import (
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
)

// PodCreate is used for creating a Pod object.
func PodCreate(client *kubernetes.Clientset, pod *corev1.Pod) (*corev1.Pod, error) {
	_, err := client.CoreV1().Pods(pod.ObjectMeta.Namespace).Create(pod)
	if errors.IsAlreadyExists(err) {
		// This will tell Kubernetes that we want this pod to be deleted immediately.
		now := int64(0)

		// Delete the Pod.
		err = client.CoreV1().Pods(pod.ObjectMeta.Namespace).Delete(pod.ObjectMeta.Name, &metav1.DeleteOptions{
			GracePeriodSeconds: &now,
		})
		if err != nil {
			return pod, err
		}

		// Create the new pod.
		_, err = client.CoreV1().Pods(pod.ObjectMeta.Namespace).Create(pod)
		if err != nil {
			return pod, err
		}
	} else if err != nil {
		return pod, err
	}

	// Wait for the pod to become available.
	limiter := time.Tick(time.Second / 10)

	for {
		pod, err = client.CoreV1().Pods(pod.ObjectMeta.Namespace).Get(pod.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			return pod, err
		}

		if pod.Status.Phase == corev1.PodRunning {
			break
		}

		<-limiter
	}

	return pod, err
}
