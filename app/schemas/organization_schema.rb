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
        json.legalName I18n.t('legal_name')
        json.address I18n.t('address')
        json.telephone I18n.t('phones')
        json.logo do
          json.set! :@type, 'ImageObject'
          json.url 'https://raw.githubusercontent.com/Hexlet/assets/master/images/hexlet_logo.svg'
        end
      end
    end
  end
end
