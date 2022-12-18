# frozen_string_literal: true

class LanguageSchema
  class << self
    include Rails.application.routes.url_helpers

    def to_builder(language, info)
      Jbuilder.new do |json|
        json.set! '@type', 'Course'
        json.name language
        json.url language_path(language)
        json.description info.description
        json.provider ProviderSchema.to_builder
      end
    end
  end
end
