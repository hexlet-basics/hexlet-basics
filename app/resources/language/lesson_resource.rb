class Language::LessonResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Lesson::Version::Info

  attributes :id, :name, :locale, :instructions

  typelize :number, nullable: true
  attribute :id do |info|
    info.version.lesson.id
  end

  typelize :number, nullable: true
  attribute :slug do |info|
    info.lesson.slug
  end

  typelize :number, nullable: true
  attribute :natural_order do |info|
    info.version.natural_order
  end
end
