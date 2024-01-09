resource "cloudflare_zone" "hexlet_basics_zone" {
  account_id = data.cloudflare_accounts.hexlet.accounts[0].id
  zone = var.domain
}

resource "cloudflare_zone" "hexlet_basics_zone_ru" {
  account_id = data.cloudflare_accounts.hexlet.accounts[0].id
  zone = var.domain_ru
}

resource "cloudflare_record" "main" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = var.domain
  value   = var.ip
  type    = "A"
  proxied = true
}

resource "cloudflare_record" "main_ru" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone_ru.id
  name    = var.domain_ru
  value   = var.ip
  type    = "A"
  proxied = true
}

resource "cloudflare_record" "ru" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = "ru.${var.domain}"
  value   = var.ip
  type    = "A"
  proxied = true
}

resource "cloudflare_record" "www" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = "www.${var.domain}"
  value   = var.domain
  type    = "CNAME"
  proxied = true
}

resource "cloudflare_record" "www-ru" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = "www.ru.${var.domain}"
  value   = "ru.${var.domain}"
  type    = "CNAME"
  proxied = true
}

resource "cloudflare_record" "bounces" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name  = "bounces.${var.domain}"
  value   = "sparkpostmail.com"
  type    = "CNAME"
  proxied = false
}

resource "cloudflare_record" "txt" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name  = "scph0819._domainkey"
  value   = "v=DKIM1; k=rsa; h=sha256; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCaY5OgnrfYY/bD07hyyqiVtk4Pxs9iQuN7u7SCNbD2d1JQyGOXcSD7t/A6VUZum6HlgOegSdi3p9gMb4wc9C6e/RQV5EblIdwABvLMYmC0CN+DDarNrF93Sejn44vjSY+kK6jEbqFBOc7qqO9k4Nep/sXb6gEsq6a9YvOHaeivFQIDAQAB"
  type    = "TXT"
}

resource "cloudflare_record" "yandex-verification" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name  = var.domain
  value   = "yandex-verification=ab7bd61d467a0773"
  type    = "TXT"
}

resource "cloudflare_record" "yandex-verification-ru" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name  = "ru.${var.domain}"
  value   = "yandex-verification=1f24f9aeb0526152"
  type    = "TXT"
}

resource "cloudflare_record" "google-verification-ru" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name  = "ru.${var.domain}"
  value   = "google-site-verification=AdJxboarIC6NOwJ9CEkIeZXdNE7DqamnPo0P7J4DJDw"
  type    = "TXT"
}

resource "cloudflare_record" "google-verification-com" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name  = var.domain
  value   = "google-site-verification=0kk3DdOciLvoog-nVFcbRzmZH65FWmNW-_aYElP0gJk"
  type    = "TXT"
}

resource "cloudflare_record" "facebook-domain-verification" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name  = var.domain
  value   = "facebook-domain-verification=d7d3em3a29yebcswwq8aa57shrc1m6"
  type    = "TXT"
}
