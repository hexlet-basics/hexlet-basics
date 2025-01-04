class Language::LessonResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Lesson::Version::Info

  attributes :id

  typelize :number, nullable: true
  attribute :slug do |info|
    info.lesson.slug
  end
end
