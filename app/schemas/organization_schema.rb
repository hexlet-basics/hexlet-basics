# frozen_string_literal: true

class OrganizationSchema
  class << self
    include Rails.application.routes.url_helpers

    def to_builder
      Jbuilder.new do |json|
        json.set! :@type, 'Organization'
        json.name 'Hexlet'
        json.url root_url
        # json.sameAs [facebook_curl,
        #              youtube_curl,
        #              twitter_curl,
        #              soundcloud_curl]
        json.legalName 'Hexlet Ltd'
        json.address 'Itälahdenkatu 22 A, 00210 Helsinki, Finland'
        json.vatID 'FI26641607'
        json.telephone I18n.t('phones')
        json.logo do
          json.set! :@type, 'ImageObject'
          json.url 'https://raw.githubusercontent.com/Hexlet/assets/master/images/hexlet_logo.svg'
        end
      end
    end
  end
end
