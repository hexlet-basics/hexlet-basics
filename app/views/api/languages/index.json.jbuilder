# frozen_string_literal: true

json.data do
  json.array! @languages do |language|
    json.attributes do
      json.call(language, :id, :name, :slug, :learn_as)
    end
    json.links do
      json.self api_language_url(language)
      json.lessons api_language_lessons_url(language, format: :json)
    end
  end
end
