# frozen_string_literal: true

class ProviderSchema
  class << self
    include Rails.application.routes.url_helpers

    def to_builder
      Jbuilder.new do |json|
        json.set! :@type, 'Organization'
        json.name 'Hexlet'
        json.url configus.https_host
        # json.sameAs [facebook_curl,
        #              youtube_curl,
        #              twitter_curl,
        #              soundcloud_curl]
        json.legalName 'Hexlet Ltd.'
        json.address 'Itälahdenkatu 22 A, 00210 Helsinki, Finland'
        json.vatID 'VAT ID: FI26641607'
        json.telephone t('phones')
      end
    end
  end
end
# json.logo do
#   json.set! :@type, 'ImageObject'
#   json.url image_url('hexlet_logo.png')
# end
