resource "github_repository_ruleset" "main" {
  repository  = github_repository.kubectl_example.name
  name        = "main"
  target      = "branch"
  enforcement = "active"

  conditions {
    ref_name {
      include = ["~DEFAULT_BRANCH"]
      exclude = []
    }
  }

  bypass_actors {
    actor_id    = 5 # repository admin role
    actor_type  = "RepositoryRole"
    bypass_mode = "always"
  }

  rules {
    pull_request {
      dismiss_stale_reviews_on_push     = true
      required_approving_review_count   = 1
      require_last_push_approval        = false
      required_review_thread_resolution = false
    }

    required_status_checks {
      required_check {
        context = "lint"
      }
      required_check {
        context = "test (1.25)"
      }
      required_check {
        context = "test (1.26)"
      }
      required_check {
        context = "e2e"
      }
      strict_required_status_checks_policy = true
    }
  }
}
