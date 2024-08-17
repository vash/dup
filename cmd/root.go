package cmd

import (
	"errors"
	"fmt"

	"dup/pkg/editor"
	"dup/pkg/util"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/completion"
)

var (
	kubeconfig string
	namespace  string
	deployment string
	edit       bool
	SilentErr  = errors.New("SilentErr")
)

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("requires exactly 1 or 2 arguments")
	}
	for _, arg := range args {
		if !util.IsValidPod(arg) {
			return fmt.Errorf("invalid argument: %s", arg)
		}
	}
	return nil
}
func NewRootCmd(ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := editor.NewEditOptions(ioStreams)
	kubeConfigFlags := defaultConfigFlags().WithWarningPrinter(o.IOStreams)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)

	var rootCmd = &cobra.Command{
		Use:               "kubectl dup [options] <resource-type, ...resource-type-n> <resource> [output-resource-prefix]",
		Short:             "Duplicate a pod out of a Deployment",
		ValidArgsFunction: completion.ResourceTypeAndNameCompletionFunc(f),
		Args:              cobra.RangeArgs(2, 3),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, args, cmd))
			cmdutil.CheckErr(o.Run())
		},
	}

	rootCmd.Flags().BoolVarP(&o.DuplicateOptions.DuplicateInnerPod, "pod", "p", false, "Duplicate pod of resource, currently only applies for: 'StatefulSet','Deployment','CronJob','Job'")
	rootCmd.Flags().BoolVarP(&o.DuplicateOptions.DisableProbes, "disable-probes", "d", true, "Disable Readiness and liveness probes for duplicated pods only (requires '-p' for complex resources)")
	rootCmd.Flags().BoolVarP(&o.DuplicateOptions.LoopCommand, "command-loop", "l", false, "Changes running command to an infinite loop (currently : \"tail -f /dev/null\"")
	rootCmd.Flags().BoolVarP(&o.SkipEdit, "skip-edit", "k", false, "Skip editing duplicated resource before creation")
	rootCmd.Flags().BoolVar(&o.WindowsLineEndings, "windows-line-endings", o.WindowsLineEndings,
		"Defaults to the line ending native to your platform.")

	kubeConfigFlags.AddFlags(rootCmd.Flags())
	matchVersionKubeConfigFlags.AddFlags(rootCmd.Flags())
	cmdutil.AddValidateFlags(rootCmd)
	o.PrintFlags.AddFlags(rootCmd)
	return rootCmd
}

func defaultConfigFlags() *genericclioptions.ConfigFlags {
	return genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
}
