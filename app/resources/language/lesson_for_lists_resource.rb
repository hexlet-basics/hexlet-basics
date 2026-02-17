class Language::LessonForListsResource < ApplicationResource
  # include Rails.application.routes.url_helpers
  typelize_from Language::Lesson::Version::Info

  # has_one :language # , resource: LanguageResource
  attributes :id, :name, :description

  typelize :string
  attribute :slug do
    it.lesson.slug
  end
end
