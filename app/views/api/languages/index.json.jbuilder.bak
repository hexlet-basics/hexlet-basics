# frozen_string_literal: true

json.data do
  json.array! @languages do |language|
    json.id language.id
    json.attributes do
      json.name language.current_version.name
      json.call(language, :slug, :learn_as)

      info = @infos_by_language[language.id]
      json.description info.description
    end
    json.links do
      json.self api_language_url(language, format: :json, locale: AppHost.locale_for_url)
    end
  end
end
