# frozen_string_literal: true

json.data do
  json.array! @lessons do |lesson|
    json.id lesson.id
    json.attributes do
      json.call(lesson, :slug)
    end
    json.links do
      json.self api_language_lesson_url(resource_language, lesson, format: :json)
    end
  end
end
