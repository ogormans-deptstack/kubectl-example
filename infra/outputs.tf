output "repository_url" {
  value       = github_repository.kubectl_example.html_url
  description = "URL of the GitHub repository"
}

output "repository_ssh_url" {
  value       = github_repository.kubectl_example.ssh_clone_url
  description = "SSH clone URL"
}
