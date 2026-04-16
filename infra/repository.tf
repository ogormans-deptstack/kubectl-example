resource "github_repository" "kubectl_example" {
  name        = var.repo_name
  description = "Generate example Kubernetes YAML manifests from your cluster's OpenAPI v3 spec"
  visibility  = "public"

  has_issues   = true
  has_projects = false
  has_wiki     = false

  allow_squash_merge = true
  allow_merge_commit = true
  allow_rebase_merge = false

  delete_branch_on_merge = true
  auto_init              = false

  homepage_url = "https://github.com/kubernetes/enhancements/issues/5571"

  topics = [
    "kubectl",
    "kubectl-plugin",
    "kubernetes",
    "openapi",
    "yaml",
    "sig-cli",
    "krew",
  ]

  lifecycle {
    prevent_destroy = true
  }
}
