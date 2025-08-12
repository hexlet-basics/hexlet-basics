class Language::Lesson::ReviewResource < ApplicationResource
  typelize_from Language::Lesson::Review

  attributes :id,
    :language_id,
    :lesson_id,
    :lesson_version_id,
    :lesson_info_id,
    :summary,
    :created_at
end
