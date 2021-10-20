# frozen_string_literal: true

json.data do
  json.array! @languages do |language|
    json.id language.id
    json.attributes do
      json.call(language, :name, :slug, :learn_as)
      json.description t("#{language}.description", scope: 'languages')
    end
    json.links do
      json.self api_language_url(language, format: :json)
    end
  end
end
