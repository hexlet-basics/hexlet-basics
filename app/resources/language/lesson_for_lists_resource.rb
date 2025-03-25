class Language::LessonForListsResource
  # include Rails.application.routes.url_helpers
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Lesson::Version::Info

  # has_one :language # , resource: LanguageResource
  attributes :id, :name, :description

  typelize :string
  attribute :slug do |info|
    info.lesson.slug
  end
end
