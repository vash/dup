package util

import (
	"fmt"
	"os"
	"regexp"

	"github.com/google/uuid"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func GetDefaultNamespace(kubeconfig string) string {
	defaultNamespace, _, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{}}).Namespace()
	if err != nil {
		fmt.Printf("Failed getting current namespace %v\n", err)
		os.Exit(1)
	}
	return defaultNamespace
}
func GetDefaultKubeconfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not get home directory:%w", err)
	}
	kubeconfigDefaultPath := homeDir + "/.kube/config"
	return kubeconfigDefaultPath
}

func GetDefaultPodName(input string) string {
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

func IsValidPod(podName string) bool {
	// Ref: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-label-names
	return isValidDNSLabel(podName)
}
func isValidDNSLabel(input string) bool {
	// Regular expression for DNS label validation as per RFC 1123
	dnsLabelRegex := regexp.MustCompile(`^[a-z0-9]([-a-z0-9]{0,61}[a-z0-9])?$`)
	return dnsLabelRegex.MatchString(input)
}
