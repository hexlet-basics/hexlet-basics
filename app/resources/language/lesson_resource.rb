class Language::LessonResource < ApplicationResource
  typelize_from Language::Lesson::Version::Info

  has_one :language, resource: LanguageResource

  attributes :id,
    :name,
    :locale,
    :instructions,
    :theory,
    :version_id,
    :description,
    :definitions,
    :tips,
    :created_at

  typelize tips: "string[]", definitions: "Array<{ name: string, description: string }>"

  typelize :number
  attribute :id do
    it.language_lesson_id
  end

  # typelize :number, nullable: true

  typelize :string, nullable: true
  attribute :prepared_code do
    it.version.prepared_code
  end

  typelize :string, nullable: true
  attribute :original_code do
    it.version.original_code
  end

  typelize :string, nullable: true
  attribute :test_code do
    it.version.test_code
  end

  typelize :number, nullable: true
  attribute :version do
    it.version.id
  end

  typelize :string
  attribute :slug do
    it.lesson.slug
  end

  typelize :number
  attribute :natural_order do
    it.version.natural_order
  end

  typelize :string, nullable: true
  attribute :source_code_url do
    ExternalLinks.lesson_source_code_curl(it.version.path_to_code, it.locale)
  end
end
