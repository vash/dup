package main

import (
	cmd "dup/cmd"
	"os"

	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func main() {
	root := cmd.NewRootCmd(genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
