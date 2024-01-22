package cmd

import (
	"errors"

	"dup/pkg/editor"

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
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}
	o.PrintFlags.AddFlags(rootCmd)
	cmdutil.AddValidateFlags(rootCmd)
	rootCmd.Flags().BoolVar(&o.WindowsLineEndings, "windows-line-endings", o.WindowsLineEndings,
		"Defaults to the line ending native to your platform.")
	rootCmd.Flags().BoolVarP(&edit, "edit", "e", false, "Edit duplicated resource before creation")
	return rootCmd
}

func defaultConfigFlags() *genericclioptions.ConfigFlags {
	return genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
}
