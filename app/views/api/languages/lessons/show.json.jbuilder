# frozen_string_literal: true

json.data do
  json.id @lesson.id
  json.attributes do
    json.call(@lesson, :slug)
    json.call(@lesson_info, :name, :theory, :description, :tips, :definitions, :instructions)
  end
  json.links do
    json.language api_language_url(resource_language, format: :json)
    json.next_lesson @lesson_version.next_lesson ? api_language_lesson_url(resource_language, @lesson_version.next_lesson, format: :json, locale: AppHost.locale_for_url) : nil
    json.prev_lesson @lesson_version.prev_lesson ? api_language_lesson_url(resource_language, @lesson_version.prev_lesson, format: :json, locale: AppHost.locale_for_url) : nil
  end
end
