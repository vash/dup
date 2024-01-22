package cmd

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

func CreateDuplicatePod(ctx context.Context, ioStreams genericclioptions.IOStreams, client *kubernetes.Clientset, deploymentName, namespace, duplicateName string, edit bool) {
	var err error
	clonedPod := clonePod(client, deploymentName, namespace, duplicateName)
	clonedPod.Spec.RestartPolicy = "Never"

	if err != nil {
		fmt.Printf("Error editing pod: %v\n", err)
		os.Exit(1)
	}
	_, err = client.CoreV1().Pods(namespace).Create(ctx, clonedPod, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error creating cloned deployment: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Cloned deployment created successfully.")
}

func clonePod(client *kubernetes.Clientset, deploymentName, namespace, duplicateName string) *corev1.Pod {
	var clonedPod corev1.Pod

	originalDeployment, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Could not find deployment: %v\n", err)
		os.Exit(1)
	}
	originalDeployment.Spec.Template.Spec.DeepCopyInto(&clonedPod.Spec)

	labelMap := make(map[string]string)
	for key, value := range originalDeployment.ObjectMeta.Labels {
		labelMap[key] = value
	}
	clonedPod.ObjectMeta.Labels = labelMap
	clonedPod.Name = duplicateName

	return &clonedPod
}
