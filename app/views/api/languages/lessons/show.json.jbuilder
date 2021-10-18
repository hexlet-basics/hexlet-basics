# frozen_string_literal: true

json.data do
  json.attributes do
    json.call(@lesson, :id, :slug)
    json.call(@lesson_info, :name, :theory, :description, :tips, :definitions, :instructions)
  end
  json.links do
    json.language api_language_url(resource_language, format: :json)
    json.next_lesson api_language_lesson_url(resource_language, @lesson_version.next_lesson, format: :json)
    json.prev_lesson api_language_lesson_url(resource_language, @lesson_version.prev_lesson, format: :json)
  end
end
