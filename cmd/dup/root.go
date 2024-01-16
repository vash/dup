package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig string
	namespace  string
	deployment string
	SilentErr  = errors.New("SilentErr")
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig file location (default is $HOME/.kube/config)")
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "deployment namespace")
	rootCmd.Flags().StringVarP(&deployment, "deployment", "d", "", "deployment name (required)")
	rootCmd.MarkFlagRequired("deployment")

	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return SilentErr
	})
}

var rootCmd = &cobra.Command{
	Use:   "dup",
	Short: "Duplicate a pod out of a Deployment",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && !IsValidPod(args[0]) {
			return fmt.Errorf("invalid pod name specified: %s", args[0])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var podName string
		if kubeconfig == "" {
			kubeconfig = getDefaultKubeconfigPath()
		}
		if namespace == "" {
			namespace = getDefaultNamespace(kubeconfig)
		}
		if len(args) == 0 {
			podName = getDefaultPodName(deployment)
		} else {
			podName = args[0]
		}
		config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
		clientset, _ := kubernetes.NewForConfig(config)

		CreateDuplicatePod(context.Background(), clientset, deployment, namespace, podName)
	},
}
