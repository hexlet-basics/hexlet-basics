class Language::OriginalLessonResource < ApplicationResource
  typelize_from Language::Lesson

  has_one :language, resource: LanguageResource

  attributes :id,
    :review,
    :slug,
    :language_id,
    :module_id,
    :created_at
end
