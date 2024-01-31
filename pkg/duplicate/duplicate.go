package duplicate

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ClonePod(client *kubernetes.Clientset, namespace, deployment, podName string) (*corev1.Pod, error) {
	var clonedPod corev1.Pod
	deploymentObj, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), deployment, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	deploymentObj.Spec.Template.Spec.DeepCopyInto(&clonedPod.Spec)

	labelMap := make(map[string]string)
	for key, value := range deploymentObj.ObjectMeta.Labels {
		labelMap[key] = value
	}
	clonedPod.ObjectMeta.Labels = labelMap
	clonedPod.Name = podName
	disableProbes(&clonedPod)

	return &clonedPod, nil
}

func disableProbes(pod *corev1.Pod) {
	containers := pod.Spec.Containers
	for i := range containers {
		containers[i].ReadinessProbe = nil
		containers[i].LivenessProbe = nil
	}
}
