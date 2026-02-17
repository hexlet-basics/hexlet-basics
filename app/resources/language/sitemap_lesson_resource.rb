class Language::SitemapLessonResource < ApplicationResource
  typelize_from Language::Lesson::Version::Info

  typelize :number
  attribute :id do
    it.language_lesson_id
  end

  typelize :string
  attribute :name do
    it.name
  end

  typelize "Locale"
  attribute :locale do
    it.locale
  end

  typelize :string
  attribute :slug do
    it.lesson.slug
  end

  typelize :number
  attribute :natural_order do
    it.version.natural_order
  end

  typelize :number
  attribute :course_id do
    it.language_id
  end
end
