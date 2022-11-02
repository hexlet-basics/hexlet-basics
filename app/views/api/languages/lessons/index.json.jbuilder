# frozen_string_literal: true

json.data do
  json.array! @lesson_infos do |lesson_info|
    json.id lesson_info.version.lesson.id
    json.attributes do
      json.call(lesson_info.version.lesson, :slug)
      json.call(lesson_info, :name, :description)
      json.natural_order lesson_info.version.natural_order
    end
    json.links do
      json.self api_language_lesson_url(resource_language, lesson_info.version.lesson, format: :json, locale: AppHost.locale_for_url)
    end
  end
end
