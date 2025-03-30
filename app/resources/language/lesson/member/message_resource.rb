class Language::Lesson::Member::MessageResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Lesson::Member::Message

  # has_one :language, resource: LanguageResource
  # has_one :language_lesson, resource: Language::LessonForListsResource

  attributes :id,
    :language_id,
    :language_lesson_id,
    :language_lesson_member_id,
    :role,
    :body,
    :created_at

  typelize :string
  attribute :language_slug do |obj|
    obj.language.slug
  end

  typelize :string
  attribute :language_lesson_slug do |obj|
    obj.language_lesson.slug
  end
end
