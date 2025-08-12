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
  attribute :slug do |review|
    review.lesson.slug
  end

  typelize :string
  attribute :language_slug do |review|
    review.language.slug
  end
end
