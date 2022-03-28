# frozen_string_literal: true

json.set! :@context, 'https://schema.org'
json.set! :@type, 'Organization'
json.name 'Hexlet'
json.url configus.https_host
# json.sameAs [facebook_curl,
#              youtube_curl,
#              twitter_curl,
#              soundcloud_curl]
json.legalName 'Hexlet Ltd.'
json.address 'UMA Esplanadi, Pohjoisesplanadi 39, 00100 Helsinki, Finland'
json.vatID 'VAT ID: FI26641607'
json.telephone t('phones')
# json.logo do
#   json.set! :@type, 'ImageObject'
#   json.url image_url('hexlet_logo.png')
# end
