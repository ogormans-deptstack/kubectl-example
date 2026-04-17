package migrate

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Status int

const (
	StatusOK Status = iota
	StatusDeprecated
	StatusRemoved
)

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusDeprecated:
		return "DEPRECATED"
	case StatusRemoved:
		return "REMOVED"
	default:
		return "UNKNOWN"
	}
}

type Manifest struct {
	APIVersion string
	Kind       string
	Group      string
	Version    string
	Name       string
}

type Result struct {
	Manifest    Manifest
	Status      Status
	Replacement string
}

func ParseManifests(data []byte) ([]Manifest, error) {
	var manifests []Manifest
	docs := bytes.Split(data, []byte("\n---"))

	for _, doc := range docs {
		trimmed := bytes.TrimSpace(doc)
		if len(trimmed) == 0 {
			continue
		}

		var raw map[string]any
		if err := yaml.Unmarshal(trimmed, &raw); err != nil {
			return nil, fmt.Errorf("parse YAML document: %w", err)
		}
		if len(raw) == 0 {
			continue
		}

		apiVersion, _ := raw["apiVersion"].(string)
		kind, _ := raw["kind"].(string)
		name := extractName(raw)
		group, version := splitAPIVersion(apiVersion)

		manifests = append(manifests, Manifest{
			APIVersion: apiVersion,
			Kind:       kind,
			Group:      group,
			Version:    version,
			Name:       name,
		})
	}

	return manifests, nil
}

func CheckAgainstAvailable(m Manifest, available map[string][]string) Result {
	versions, groupExists := available[m.Group]
	if !groupExists {
		return Result{
			Manifest:    m,
			Status:      StatusRemoved,
			Replacement: "",
		}
	}

	for _, v := range versions {
		if v == m.Version {
			return Result{Manifest: m, Status: StatusOK}
		}
	}

	if len(versions) > 0 {
		replacement := m.Group + "/" + versions[len(versions)-1]
		if m.Group == "" {
			replacement = versions[len(versions)-1]
		}
		return Result{
			Manifest:    m,
			Status:      StatusDeprecated,
			Replacement: replacement,
		}
	}

	return Result{
		Manifest: m,
		Status:   StatusRemoved,
	}
}

func AnalyzeBytes(data []byte, available map[string][]string) ([]Result, error) {
	manifests, err := ParseManifests(data)
	if err != nil {
		return nil, err
	}

	results := make([]Result, len(manifests))
	for i, m := range manifests {
		results[i] = CheckAgainstAvailable(m, available)
	}
	return results, nil
}

func splitAPIVersion(apiVersion string) (group, version string) {
	if i := strings.LastIndex(apiVersion, "/"); i >= 0 {
		return apiVersion[:i], apiVersion[i+1:]
	}
	return "", apiVersion
}

func extractName(raw map[string]any) string {
	meta, ok := raw["metadata"].(map[string]any)
	if !ok {
		return ""
	}
	name, _ := meta["name"].(string)
	return name
}
