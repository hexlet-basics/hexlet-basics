# frozen_string_literal: true

json.data do
  json.array! @lessons do |lesson|
    json.attributes do
      json.call(lesson, :id, :slug)
    end
    json.links do
      json.self api_language_lesson_url(resource_language, lesson, format: :json)
    end
  end
end
