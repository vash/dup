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
	// Validate the number of arguments
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("requires exactly 1 or 2 arguments")
	}

	// Validate each argument using ValidName function
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
		Use:               "dup",
		Short:             "Duplicate a pod out of a Deployment",
		ValidArgsFunction: completion.ResourceTypeAndNameCompletionFunc(f),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, args, cmd))
			cmdutil.CheckErr(o.Run())
		},
	}
	o.PrintFlags.AddFlags(rootCmd)
	cmdutil.AddValidateFlags(rootCmd)
	rootCmd.Flags().BoolVar(&o.WindowsLineEndings, "windows-line-endings", o.WindowsLineEndings,
		"Defaults to the line ending native to your platform.")
	rootCmd.Flags().BoolVarP(&o.SkipEdit, "skip edit", "s", false, "Skip editing duplicated resource before creation")
	return rootCmd
}

func defaultConfigFlags() *genericclioptions.ConfigFlags {
	return genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
}
