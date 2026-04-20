package migrate

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ogormans-deptstack/kubectl-schemagen/internal/cli"
	"github.com/ogormans-deptstack/kubectl-schemagen/pkg/migrate"
)

func NewCommand() *cobra.Command {
	var kubeconfig string

	cmd := &cobra.Command{
		Use:   "migrate [FILE...]",
		Short: "Detect deprecated or removed Kubernetes APIs in manifests",
		Long: `Reads YAML manifests and detects deprecated or removed API versions by
comparing against the connected cluster's OpenAPI schema. Reports which
resources use outdated API versions and suggests replacements.`,
		Example: `  kubectl schemagen migrate deployment.yaml
  kubectl schemagen migrate manifests/*.yaml`,
		Aliases: []string{"mig"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrate(args, kubeconfig)
		},
	}

	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig")

	return cmd
}

func runMigrate(files []string, kubeconfig string) error {
	available, err := cli.LoadAvailableAPIs(kubeconfig)
	if err != nil {
		return err
	}

	hasIssues := false

	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		results, err := migrate.AnalyzeBytes(data, available)
		if err != nil {
			return fmt.Errorf("analyze %s: %w", path, err)
		}

		for _, r := range results {
			switch r.Status {
			case migrate.StatusOK:
				fmt.Printf("  OK  %s %s/%s\n", path, r.Manifest.APIVersion, r.Manifest.Kind)
			case migrate.StatusDeprecated:
				fmt.Printf(" DEP  %s %s/%s -> %s\n", path, r.Manifest.APIVersion, r.Manifest.Kind, r.Replacement)
				hasIssues = true
			case migrate.StatusRemoved:
				if r.Replacement != "" {
					fmt.Printf(" REM  %s %s/%s -> %s\n", path, r.Manifest.APIVersion, r.Manifest.Kind, r.Replacement)
				} else {
					fmt.Printf(" REM  %s %s/%s (no replacement found)\n", path, r.Manifest.APIVersion, r.Manifest.Kind)
				}
				hasIssues = true
			}
		}
	}

	if hasIssues {
		return fmt.Errorf("deprecated or removed APIs found")
	}
	return nil
}
