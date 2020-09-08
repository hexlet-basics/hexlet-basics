resource "google_cloudbuild_trigger" "app" {
  project     = google_project.hexlet_basics.project_id
  description = "app"

  trigger_template {
    branch_name = "master"
    repo_name   = var.repositories["hexlet_basics"]
  }

  filename = "services/app/cloudbuild.yaml"

  substitutions = {
    _SLACK_WEBHOOK = var.slack_codebuild_webhook
  }
}

resource "google_cloudbuild_trigger" "app-tag" {
  project     = var.project_name
  description = "build with tag"

  trigger_template {
    tag_name  = "v\\d.*"
    repo_name = var.repositories["hexlet_basics"]
  }

  filename = "services/app/cloudbuild.yaml"

  substitutions = {
    _SLACK_WEBHOOK = var.slack_codebuild_webhook
  }
}

resource "google_cloudbuild_trigger" "caddy" {
  project     = var.project_name
  description = "caddy"

  trigger_template {
    tag_name  = "v\\d.*"
    # project   = var.project_name
    repo_name = var.repositories["hexlet_basics"]
  }

  substitutions = {
    _SLACK_WEBHOOK = var.slack_codebuild_webhook
  }

  filename = "services/caddy/cloudbuild.yaml"
}

# resource "google_cloudbuild_trigger" "exercises_php" {
#   project     = var.project_name
#   description = "exercises_php"


#   trigger_template {
#     branch_name = "master"
#     project     = var.project_name
#     repo_name   = var.repositories["exercises_php"]
#   }


#   substitutions = {
#     _SLACK_WEBHOOK = var.slack_codebuild_webhook
#   }


#   filename = "cloudbuild.yaml"
# }


# resource "google_cloudbuild_trigger" "exercises_javascript" {
#   project     = var.project_name
#   description = "exercises_javascript"


#   trigger_template {
#     branch_name = "master"
#     project     = var.project_name
#     repo_name   = var.repositories["exercises_javascript"]
#   }


#   substitutions = {
#     _SLACK_WEBHOOK = var.slack_codebuild_webhook
#   }


#   filename = "cloudbuild.yaml"
# }


# resource "google_cloudbuild_trigger" "exercises_java" {
#   project     = var.project_name
#   description = "exercises_java"


#   trigger_template {
#     branch_name = "master"
#     project     = var.project_name
#     repo_name   = var.repositories["exercises_java"]
#   }


#   substitutions = {
#     _SLACK_WEBHOOK = var.slack_codebuild_webhook
#   }


#   filename = "cloudbuild.yaml"
# }


# resource "google_cloudbuild_trigger" "exercises_python" {
#   project     = var.project_name
#   description = "exercises_python"


#   trigger_template {
#     branch_name = "master"
#     project     = var.project_name
#     repo_name   = var.repositories["exercises_python"]
#   }


#   substitutions = {
#     _SLACK_WEBHOOK = var.slack_codebuild_webhook
#   }


#   filename = "cloudbuild.yaml"
# }

