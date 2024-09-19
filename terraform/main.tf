# --------------------------------------
# SPACES
# --------------------------------------
resource "digitalocean_spaces_bucket" "sitemap_bucket" {
  name   = var.do_spaces_sitemap_bucket
  region = var.do_spaces_region
  acl    = "private"
  lifecycle_rule {
    enabled = true
    expiration {
      days = 30
    }
  }
}

# --------------------------------------
# PROJECT
# --------------------------------------

resource "digitalocean_project" "hexlet_basics_project" {
  name        = "Hexlet Basics"
  description = "A project to represent Hexlet Basics resources."
  purpose     = "Web Application"
  environment = "Production"
  resources   = [
  ]
}
