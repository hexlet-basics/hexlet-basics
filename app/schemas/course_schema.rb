# frozen_string_literal: true

class CourseSchema
  class << self
    include Rails.application.routes.url_helpers

    def to_builder(language, info)
      Jbuilder.new do |json|
        json.set! '@type', 'Course'
        json.name info.header
        json.url language_url(language.slug, locale: AppHost.locale_for_url(info.locale))
        json.description info.description
        json.provider ProviderSchema.to_builder
      end
    end
  end
end
