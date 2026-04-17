package migrate

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ogormans-deptstack/kubectl-generate/internal/cli"
	"github.com/ogormans-deptstack/kubectl-generate/pkg/migrate"
	"github.com/ogormans-deptstack/kubectl-generate/pkg/openapi"
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
	doc, err := cli.LoadClusterDoc(kubeconfig)
	if err != nil {
		return err
	}

	available := buildAvailableAPIs(doc)
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
				fmt.Printf(" REM  %s %s/%s (no replacement found)\n", path, r.Manifest.APIVersion, r.Manifest.Kind)
				hasIssues = true
			}
		}
	}

	if hasIssues {
		return fmt.Errorf("deprecated or removed APIs found")
	}
	return nil
}

func buildAvailableAPIs(doc *openapi.Document) map[string][]string {
	available := make(map[string][]string)
	schemas := doc.ComponentSchemas()
	if schemas == nil {
		return available
	}

	seen := make(map[string]map[string]bool)
	for _, schema := range schemas {
		schemaMap, ok := schema.(map[string]any)
		if !ok {
			continue
		}
		ext, ok := schemaMap["x-kubernetes-group-version-kind"]
		if !ok {
			continue
		}
		arr, ok := ext.([]any)
		if !ok {
			continue
		}
		for _, item := range arr {
			m, ok := item.(map[string]any)
			if !ok {
				continue
			}
			group, _ := m["group"].(string)
			version, _ := m["version"].(string)
			if seen[group] == nil {
				seen[group] = make(map[string]bool)
			}
			if !seen[group][version] {
				seen[group][version] = true
				available[group] = append(available[group], version)
			}
		}
	}
	return available
}
