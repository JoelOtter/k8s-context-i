package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JoelOtter/k8s-context-i/internal/k8s"
	"github.com/JoelOtter/k8s-context-i/internal/ui"
	"github.com/spf13/cobra"
)

func main() {
	var debug bool

	cmd := &cobra.Command{
		Use: "k8s-context-i",
		RunE: func(cmd *cobra.Command, args []string) error {
			contexts, err := k8s.GetContexts()
			if err != nil {
				return err
			}
			if err := ui.ShowUI(contexts); err != nil {
				return fmt.Errorf("failed to show UI: %w", err)
			}
			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.PersistentFlags().BoolVar(
		&debug,
		"debug",
		false,
		"Show debug output",
	)

	if err := cmd.Execute(); err != nil {
		if debug {
			log.Fatalln(err)
		}
		os.Exit(1)
	}
}
