# frozen_string_literal: true

json.data do
  json.attributes do
    json.call(@language, :id, :name, :slug, :learn_as)
  end
  json.links do
    json.first_lesson api_language_lesson_url(@language, @language.lessons.first, format: :json)
  end
  json.relationships do
    json.lessons do
      json.array! @language.lessons do |lesson|
        json.attributes do
          json.call(lesson, :id, :slug)
        end
        json.links do
          json.self api_language_lesson_url(@language, lesson, format: :json)
        end
      end
    end
  end
end
