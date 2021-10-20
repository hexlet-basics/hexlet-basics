# frozen_string_literal: true

json.data do
  @current_module_versions.each do |module_version|
    json.array! module_version.lesson_versions.each do |lesson_version|
      lesson_info = @infos_by_lesson[lesson_version.id]

      json.id lesson_info.version.lesson.id
      json.attributes do
        json.call(lesson_info.version.lesson, :slug)
        json.call(lesson_info, :name, :description)
        json.natural_order lesson_version.natural_order
      end
      json.links do
        json.self api_language_lesson_url(resource_language, lesson_info.version.lesson, format: :json)
      end
    end
  end
end
