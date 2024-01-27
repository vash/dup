package duplicate

import (
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
)

func ClonePods(client *kubernetes.Clientset, deployments []*resource.Info, namespace, duplicateName string) []*resource.Info {
	var ret []*resource.Info
	for _, deployment := range deployments {
		ret = append(ret, &(resource.Info{Object: ClonePod(client, deployment.Object, namespace, duplicateName)}))
	}
	return ret
}
func ClonePod(client *kubernetes.Clientset, deployment runtime.Object, namespace, duplicateName string) *corev1.Pod {
	var clonedPod corev1.Pod
	var originalDeployment appsv1.Deployment

	unstructuredDeployment := deployment.(*unstructured.Unstructured).UnstructuredContent()
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredDeployment, &originalDeployment)
	if err != nil {
		fmt.Printf("Failed to convert deployment: %s", err)
		os.Exit(1)
	}
	originalDeployment.Spec.Template.Spec.DeepCopyInto(&clonedPod.Spec)

	labelMap := make(map[string]string)
	for key, value := range originalDeployment.ObjectMeta.Labels {
		labelMap[key] = value
	}
	clonedPod.ObjectMeta.Labels = labelMap
	clonedPod.Name = duplicateName
	clonedPod.Spec.RestartPolicy = "Never"
	return &clonedPod
}
