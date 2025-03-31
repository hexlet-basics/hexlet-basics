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

  typelize :number
  attribute :user_id do |obj|
    obj.language_lesson_member.user_id
  end

  typelize :string
  attribute :language_lesson_slug do |obj|
    obj.language_lesson.slug
  end

  typelize :string
  attribute :language_lesson_name do |obj|
    info = obj.language_lesson.localed_info
    info.name
  end
end
