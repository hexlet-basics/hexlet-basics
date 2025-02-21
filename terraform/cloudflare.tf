# TODO: в аккаунте есть ещё 5 dns записей, которые видимо добавляли руками
# проверить нужны ли они, и если да то импортировать в terraform

# --------------------------------------
# ACCOUNTS
# --------------------------------------
data "cloudflare_accounts" "hexlet" {
  name = "Hexlet Production"
}

# --------------------------------------
# code-basics.ru
# --------------------------------------
resource "cloudflare_zone" "hexlet_basics_zone_ru" {
  account = {
    id = data.cloudflare_accounts.hexlet.result[0].id
  }
  name = var.domain_ru
}

resource "cloudflare_dns_record" "main_ru_1" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone_ru.id
  name    = var.domain_ru
  content = var.k8s_data.ip1
  type    = "A"
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "main_ru_www" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone_ru.id
  name    = "www.${var.domain_ru}"
  content = var.domain_ru
  type    = "CNAME"
  proxied = true
   ttl    = 1
}

resource "cloudflare_ruleset" "http_request_dynamic_redirect_main_ru" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone_ru.id
  name    = "default"
  kind    = "zone"
  phase   = "http_request_dynamic_redirect"

  rules = [{
    action      = "redirect"
    description = "Redirect from ru domain to main domain"
    enabled     = true
    expression = "true"

    action_parameters = {
      from_value = {
        preserve_query_string = true
        status_code           = 301
        target_url = {
          value = "https://${var.domain}/ru"
        }
      }
    }
  }]
}

# --------------------------------------
# code-basics.com
# --------------------------------------
resource "cloudflare_zone" "hexlet_basics_zone" {
  account = {
    id = data.cloudflare_accounts.hexlet.result[0].id
  }
  name = var.domain
}

resource "cloudflare_dns_record" "main_1" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = var.domain
  content = var.k8s_data.ip1
  type    = "A"
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "main_2" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = var.domain
  content = var.k8s_data.ip2
  type    = "A"
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "main_3" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = var.domain
  content = var.k8s_data.ip3
  type    = "A"
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "ru" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = "ru.${var.domain}"
  content = var.domain
  type    = "CNAME"
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "www" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = "www.${var.domain}"
  content = var.domain
  type    = "CNAME"
  proxied = true
  ttl     = 1
}

# NOTE: не работает так как сертификат не покрывает поддомен 4-го уровня
resource "cloudflare_dns_record" "www_ru" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = "www.ru.${var.domain}"
  content = "ru.${var.domain}"
  type    = "CNAME"
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "bounces" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = "bounces.${var.domain}"
  content = "sparkpostmail.com"
  type    = "CNAME"
  proxied = false
  ttl     = 1
}

resource "cloudflare_dns_record" "txt" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = "scph0819._domainkey.${var.domain}"
  content = "v=DKIM1; k=rsa; h=sha256; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCaY5OgnrfYY/bD07hyyqiVtk4Pxs9iQuN7u7SCNbD2d1JQyGOXcSD7t/A6VUZum6HlgOegSdi3p9gMb4wc9C6e/RQV5EblIdwABvLMYmC0CN+DDarNrF93Sejn44vjSY+kK6jEbqFBOc7qqO9k4Nep/sXb6gEsq6a9YvOHaeivFQIDAQAB"
  type    = "TXT"
  ttl     = 1
}

resource "cloudflare_dns_record" "yandex_verification" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = var.domain
  content = "yandex-verification=ab7bd61d467a0773"
  type    = "TXT"
  ttl     = 1
}

resource "cloudflare_dns_record" "yandex_verification_ru" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = "ru.${var.domain}"
  content = "yandex-verification=1f24f9aeb0526152"
  type    = "TXT"
  ttl     = 1
}

resource "cloudflare_dns_record" "google_verification_ru" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = "ru.${var.domain}"
  content = "google-site-verification=AdJxboarIC6NOwJ9CEkIeZXdNE7DqamnPo0P7J4DJDw"
  type    = "TXT"
  ttl     = 1
}

resource "cloudflare_dns_record" "google_verification_com" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = var.domain
  content = "google-site-verification=0kk3DdOciLvoog-nVFcbRzmZH65FWmNW-_aYElP0gJk"
  type    = "TXT"
  ttl     = 1
}

resource "cloudflare_dns_record" "facebook_domain_verification" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = var.domain
  content = "facebook-domain-verification=d7d3em3a29yebcswwq8aa57shrc1m6"
  type    = "TXT"
  ttl     = 1
}

resource "cloudflare_ruleset" "http_request_dynamic_redirect_main" {
  zone_id = resource.cloudflare_zone.hexlet_basics_zone.id
  name    = "default"
  kind    = "zone"
  phase   = "http_request_dynamic_redirect"

  rules = [{
    action      = "redirect"
    description = "Redirect from ru subdomain to main domain"
    enabled     = true
    expression = "(http.host eq \"ru.${var.domain}\")"

    action_parameters = {
      from_value = {
        preserve_query_string = true
        status_code           = 301
        target_url = {
          value = "https://${var.domain}/ru"
        }
      }
    }
  }]
}
