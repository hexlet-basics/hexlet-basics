# frozen_string_literal: true

json.data do
  json.id @language.id
  json.attributes do
    json.name @language.current_version.name
    json.call(@language, :slug, :learn_as)
    json.description @language_info.description
  end
  json.links do
    json.first_lesson api_language_lesson_url(@language, @language.lessons.first, format: :json, locale: AppHost.locale_for_url)
  end
  json.relationships do
    json.lessons do
      json.array! @language.lessons do |lesson|
        json.id lesson.id
        json.attributes do
          json.call(lesson, :slug)
        end
        json.links do
          json.self api_language_lesson_url(@language, lesson, format: :json, locale: AppHost.locale_for_url)
        end
      end
    end
  end
end
