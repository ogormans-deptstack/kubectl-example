# kubectl-example -- Agent Context

## Project

kubectl plugin that generates example YAML from OpenAPI v3 specs. Go module at `github.com/ogormans-deptstack/kubectl-example`. Apache 2.0.

## Conventions

- Go 1.25+, strict linting via golangci-lint v2
- Table-driven tests, e2e against kind cluster with server-side dry-run
- Factory pattern: `ensureXXXCRD()` helpers for CRD test groups
- No `as any` / `@ts-ignore` equivalents -- no lint suppression
- Commit messages: no `Fixes #N` or auto-close keywords (Prow flags these)
- Use SSH for git push: `git@github.com:ogormans-deptstack/kubectl-example.git`

## Key Paths

| What | Where |
|------|-------|
| CLI entry | `cmd/kubectl-example/main.go` |
| Generator | `pkg/generator/openapi_generator.go` |
| YAML emitter | `pkg/generator/yaml.go` |
| Schema fetcher | `pkg/openapi/fetcher.go` |
| Defaults | `pkg/defaults/defaults.go` |
| Flags | `pkg/flags/flags.go` |
| E2e tests | `test/e2e/e2e_test.go` |
| CI | `.github/workflows/ci.yml` |
| CronTab CRD fixture | `test/fixtures/crontab-crd.yaml` |

## sig-cli Engagement

- **KEP**: kubernetes/enhancements#5576 (PR), tracking issue #5571
- **Meeting**: Presented at sig-cli bi-weekly, April 2026. Positive reception.
- **Contacts**: ardaguclu (sig-cli member, invited to meeting), soltysh (sig-cli, expressed interest), eddiezane (sig-cli lead, KEP approver)
- **Target**: v1.37 alpha (`kubectl alpha example`), earliest Enhancements Freeze ~May 2026

## Banked: Issue #5571 Response Draft

**Target posting date: week of April 21-25, 2026**

This response needs refinement before posting -- update test counts, CRD coverage, and open questions based on meeting feedback. The meeting already happened (April 16), so rewrite the intro accordingly and remove the "can I join" phrasing. Retarget timeline question from 1.36 to 1.37.

```markdown
/reopen
/remove-lifecycle rotten

Following up from the sig-cli bi-weekly discussion on April 16 -- thanks @ardaguclu for the invite.

Working prototype is at https://github.com/ogormans-deptstack/kubectl-example. It reads the cluster's OpenAPI v3 spec and generates apply-ready YAML for any resource type, including CRDs.

**How it works:**

The plugin connects to the cluster via the discovery API, fetches the OpenAPI v3 schemas for all group-versions, and walks the schema tree to produce a minimal valid manifest. Required fields are always included, and important optional fields (strategy, ports, resources, selectors) are pulled in based on heuristics. Labels, selectors, and template metadata are wired up consistently.

**Demo output:**

$ kubectl example Deployment --name=web --image=myapp:v2 --replicas=5
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
  labels:
    app.kubernetes.io/name: web
spec:
  replicas: 5
  selector:
    matchLabels:
      app.kubernetes.io/name: web
  strategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app.kubernetes.io/name: web
    spec:
      containers:
        - name: web
          image: "myapp:v2"
          ports:
            - containerPort: 80
          resources:
            limits:
              cpu: "500m"
              memory: "256Mi"
            requests:
              cpu: "250m"
              memory: "128Mi"

**Current state:**

- 13 core resource types pass server-side dry-run validation
- CRD support: CronTab, 10 Gateway API CRDs, [UPDATE: add Argo/Crossplane/cert-manager counts]
- Override flags: --name, --image, --replicas, --set key=value
- [UPDATE: test counts after CRD expansion]
- CI green, distributed via krew [UPDATE: once krew manifest lands]

**Next steps from the meeting discussion:**

[UPDATE: incorporate actual meeting feedback here]

1. [UPDATE: specific action items from meeting]
2. Targeting v1.37 for kubectl alpha example -- KEP #5576 needs updating to reflect this
3. [UPDATE: any design changes agreed at meeting]

cc @soltysh
```

**Before posting, update:**
- [ ] Test counts (unit + e2e) after CRD expansion
- [ ] CRD coverage list (add Argo, Crossplane, cert-manager if done)
- [ ] Meeting feedback items as concrete next steps
- [ ] Krew status (if manifest is ready by then)
- [ ] Remove placeholder [UPDATE] markers
- [ ] Verify CI is still green

## Release Plan

- GoReleaser for cross-platform builds (darwin/linux/windows, amd64/arm64)
- Krew manifest in repo root
- GitHub Actions release workflow triggered on tag push (v*)
- First release: v0.1.0

## CRD Test Expansion

| CRD Group | Install Method | Key Kinds | Status |
|-----------|---------------|-----------|--------|
| CronTab | local fixture | CronTab | done |
| Gateway API | remote URL (v1.2.1) | HTTPRoute, Gateway, +8 | done |
| Argo Workflows | remote URL | Workflow, CronWorkflow, WorkflowTemplate | planned |
| Crossplane | remote URL | Composition, CompositeResourceDefinition, Provider | planned |
| cert-manager | remote URL | Certificate, Issuer, ClusterIssuer | planned |
