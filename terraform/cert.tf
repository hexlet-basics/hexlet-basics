resource "yandex_cm_certificate" "code_basics_com" {
  name    = "code-basics-com"
  domains = ["code-basics.com"]
  deletion_protection = true

  managed {
    challenge_type  = "DNS_CNAME"
  }
}
