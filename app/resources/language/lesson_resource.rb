class Language::LessonResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Lesson::Version::Info

  attributes :id,
    :name,
    :locale,
    :instructions,
    :theory,
    :version_id,
    :description,
    :definitions,
    :tips

  typelize tips: "string[]", definitions: "Array<{ name: string, description: string }>"

  typelize :number
  attribute :id do |info|
    info.language_lesson_id
  end

  # typelize :number, nullable: true

  typelize :string, nullable: true
  attribute :prepared_code do |info|
    info.version.prepared_code
  end

  typelize :string, nullable: true
  attribute :original_code do |info|
    info.version.original_code
  end

  typelize :string, nullable: true
  attribute :test_code do |info|
    info.version.test_code
  end

  typelize :number, nullable: true
  attribute :version do |info|
    info.version.id
  end

  typelize :string
  attribute :slug do |info|
    info.lesson.slug
  end

  typelize :number
  attribute :natural_order do |info|
    info.version.natural_order
  end

  typelize :string, nullable: true
  attribute :source_code_url do |info|
    ExternalLinks.lesson_source_code_curl(info)
  end
end
