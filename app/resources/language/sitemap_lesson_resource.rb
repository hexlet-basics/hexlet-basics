class Language::SitemapLessonResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Lesson::Version::Info

  typelize :number
  attribute :id do |info|
    info.language_lesson_id
  end

  typelize :string
  attribute :name do |info|
    info.name
  end

  typelize "Locale"
  attribute :locale do |info|
    info.locale
  end

  typelize :string
  attribute :slug do |info|
    info.lesson.slug
  end

  typelize :number
  attribute :natural_order do |info|
    info.version.natural_order
  end

  typelize :number
  attribute :course_id do |info|
    info.language_id
  end
end
