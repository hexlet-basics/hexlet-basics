class Language::Lesson::ReviewResource < ApplicationResource
  typelize_from Language::Lesson::Review

  attributes :id,
    :locale,
    :language_id,
    :language_lesson_id,
    :language_lesson_version_id,
    :language_lesson_version_info_id,
    :summary,
    :created_at

  typelize :string
  attribute :slug do
    it.lesson.slug
  end

  typelize :number
  attribute :lesson_natural_order do
    it.lesson.natural_order
  end

  typelize :string
  attribute :language_slug do
    it.language.slug
  end
end
