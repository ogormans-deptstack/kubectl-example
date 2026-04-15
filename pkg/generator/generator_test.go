package generator

import (
	"bytes"
	"strings"
	"testing"
)

func TestResourceGeneratorInterface(t *testing.T) {
	t.Run("SupportedTypes returns at least 13 core types", func(t *testing.T) {
		gen := newTestGenerator(t)
		types := gen.SupportedTypes()
		if len(types) < 13 {
			t.Errorf("expected at least 13 supported types, got %d: %v", len(types), types)
		}

		coreTypes := []string{
			"Pod", "Deployment", "Service", "ConfigMap", "Secret",
			"Job", "CronJob", "Ingress", "NetworkPolicy",
			"StatefulSet", "DaemonSet", "PersistentVolumeClaim",
			"HorizontalPodAutoscaler",
		}
		supported := make(map[string]bool)
		for _, t := range types {
			supported[t] = true
		}
		for _, ct := range coreTypes {
			if !supported[ct] {
				t.Errorf("core type %s not in SupportedTypes()", ct)
			}
		}
	})
}

func TestGenerateYAML(t *testing.T) {
	coreTypes := []struct {
		name             string
		requiredFields   []string
		overrides        map[string]string
		expectedContains []string
	}{
		{
			name:           "Pod",
			requiredFields: []string{"apiVersion: v1", "kind: Pod", "metadata:", "spec:", "containers:"},
		},
		{
			name:             "Deployment",
			requiredFields:   []string{"apiVersion: apps/v1", "kind: Deployment", "spec:", "replicas:", "selector:", "template:"},
			overrides:        map[string]string{"replicas": "5", "name": "web"},
			expectedContains: []string{"replicas: 5", "name: web"},
		},
		{
			name:           "Service",
			requiredFields: []string{"apiVersion: v1", "kind: Service", "spec:", "ports:"},
		},
		{
			name:           "ConfigMap",
			requiredFields: []string{"apiVersion: v1", "kind: ConfigMap", "metadata:"},
		},
		{
			name:           "Secret",
			requiredFields: []string{"apiVersion: v1", "kind: Secret", "metadata:"},
		},
		{
			name:           "Job",
			requiredFields: []string{"apiVersion: batch/v1", "kind: Job", "spec:", "template:"},
		},
		{
			name:           "CronJob",
			requiredFields: []string{"apiVersion: batch/v1", "kind: CronJob", "spec:", "schedule:"},
		},
		{
			name:           "Ingress",
			requiredFields: []string{"apiVersion: networking.k8s.io/v1", "kind: Ingress", "spec:", "rules:"},
		},
		{
			name:           "NetworkPolicy",
			requiredFields: []string{"apiVersion: networking.k8s.io/v1", "kind: NetworkPolicy", "spec:", "podSelector:"},
		},
		{
			name:           "StatefulSet",
			requiredFields: []string{"apiVersion: apps/v1", "kind: StatefulSet", "spec:", "serviceName:"},
		},
		{
			name:           "DaemonSet",
			requiredFields: []string{"apiVersion: apps/v1", "kind: DaemonSet", "spec:", "selector:", "template:"},
		},
		{
			name:           "PersistentVolumeClaim",
			requiredFields: []string{"apiVersion: v1", "kind: PersistentVolumeClaim", "spec:", "accessModes:", "resources:"},
		},
		{
			name:           "HorizontalPodAutoscaler",
			requiredFields: []string{"apiVersion: autoscaling/v2", "kind: HorizontalPodAutoscaler", "spec:", "scaleTargetRef:"},
		},
	}

	gen := newTestGenerator(t)

	for _, tc := range coreTypes {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("generates valid YAML with required fields", func(t *testing.T) {
				var buf bytes.Buffer
				overrides := tc.overrides
				if overrides == nil {
					overrides = map[string]string{}
				}
				err := gen.Generate(tc.name, overrides, &buf)
				if err != nil {
					t.Fatalf("Generate(%s) failed: %v", tc.name, err)
				}
				yaml := buf.String()
				for _, field := range tc.requiredFields {
					if !strings.Contains(yaml, field) {
						t.Errorf("YAML for %s missing required field %q\ngot:\n%s", tc.name, field, yaml)
					}
				}
			})

			if tc.overrides != nil {
				t.Run("respects overrides", func(t *testing.T) {
					var buf bytes.Buffer
					err := gen.Generate(tc.name, tc.overrides, &buf)
					if err != nil {
						t.Fatalf("Generate(%s) with overrides failed: %v", tc.name, err)
					}
					yaml := buf.String()
					for _, expected := range tc.expectedContains {
						if !strings.Contains(yaml, expected) {
							t.Errorf("YAML for %s with overrides missing %q\ngot:\n%s", tc.name, expected, yaml)
						}
					}
				})
			}
		})
	}
}

func newTestGenerator(t *testing.T) ResourceGenerator {
	t.Helper()
	t.Skip("ResourceGenerator implementation not yet available")
	return nil
}
