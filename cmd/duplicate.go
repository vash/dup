package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/google/uuid"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
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

func getDefaultNamespace(kubeconfig string) string {
	defaultNamespace, _, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{}}).Namespace()
	if err != nil {
		fmt.Printf("Failed getting current namespace %v\n", err)
		os.Exit(1)
	}
	return defaultNamespace
}
func getDefaultKubeconfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not get home directory:%w", err)
	}
	kubeconfigDefaultPath := homeDir + "/.kube/config"
	return kubeconfigDefaultPath
}

func getDefaultPodName(input string) string {
	var result string
	var suffix = "-dup-" + uuid.New().String()[:4]
	var maxPodNameLength = 63 - len(suffix) //RFC1035

	if len(input) > (maxPodNameLength) {
		result = input[:maxPodNameLength] + suffix
	} else {
		result = input + suffix
	}
	return result
}

func isValidPod(podName string) bool {
	// Ref: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-label-names
	return isValidDNSLabel(podName)
}
func isValidDNSLabel(input string) bool {
	// Regular expression for DNS label validation as per RFC 1123
	dnsLabelRegex := regexp.MustCompile(`^[a-z0-9]([-a-z0-9]{0,61}[a-z0-9])?$`)
	return dnsLabelRegex.MatchString(input)
}
