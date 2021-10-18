variable "domain" {
  default = "code-basics.com"
}

variable "ip" {
  default = "157.245.27.205"
}

resource "cloudflare_record" "main" {
  domain  = var.domain
  name    = var.domain
  value   = var.ip
  type    = "A"
  proxied = true
}

resource "cloudflare_record" "main-ru" {
  domain  = "code-basics.ru"
  name    = "code-basics.ru"
  value   = var.ip
  type    = "A"
  proxied = true
}

resource "cloudflare_record" "ru" {
  domain  = var.domain
  name    = "ru.${var.domain}"
  value   = var.ip
  type    = "A"
  proxied = true
}

resource "cloudflare_record" "www" {
  domain  = var.domain
  name    = "www.${var.domain}"
  value   = var.domain
  type    = "CNAME"
  proxied = true
}

resource "cloudflare_record" "www-ru" {
  domain  = var.domain
  name    = "www.ru.${var.domain}"
  value   = "ru.${var.domain}"
  type    = "CNAME"
  proxied = true
}

resource "cloudflare_record" "bounces" {
  domain  = var.domain
  name  = "bounces"
  value   = "sparkpostmail.com"
  type    = "CNAME"
  proxied = false
}

resource "cloudflare_record" "txt" {
  domain  = var.domain
  name  = "scph0819._domainkey"
  value   = "v=DKIM1; k=rsa; h=sha256; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCaY5OgnrfYY/bD07hyyqiVtk4Pxs9iQuN7u7SCNbD2d1JQyGOXcSD7t/A6VUZum6HlgOegSdi3p9gMb4wc9C6e/RQV5EblIdwABvLMYmC0CN+DDarNrF93Sejn44vjSY+kK6jEbqFBOc7qqO9k4Nep/sXb6gEsq6a9YvOHaeivFQIDAQAB"
  type    = "TXT"
}

resource "cloudflare_record" "yandex-verification" {
  domain  = var.domain
  name  = var.domain
  value   = "yandex-verification=a1675c88bd4b0ade"
  type    = "TXT"
}

resource "cloudflare_record" "yandex-verification-ru" {
  domain  = var.domain
  name  = "ru.${var.domain}"
  value   = "yandex-verification=a1675c88bd4b0ade"
  type    = "TXT"
}

resource "cloudflare_record" "google-verification-ru" {
  domain  = var.domain
  name  = "ru.${var.domain}"
  value   = "google-site-verification=AdJxboarIC6NOwJ9CEkIeZXdNE7DqamnPo0P7J4DJDw"
  type    = "TXT"
}

resource "cloudflare_record" "google-verification-com" {
  domain  = var.domain
  name  = var.domain
  value   = "google-site-verification=0kk3DdOciLvoog-nVFcbRzmZH65FWmNW-_aYElP0gJk"
  type    = "TXT"
}

resource "cloudflare_record" "facebook-domain-verification" {
  domain  = var.domain
  name  = var.domain
  value   = "facebook-domain-verification=d7d3em3a29yebcswwq8aa57shrc1m6"
  type    = "TXT"
}

