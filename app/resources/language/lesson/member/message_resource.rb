class Language::Lesson::Member::MessageResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Language::Lesson::Member::Message

  attributes :id,
    :language_id,
    :language_lesson_id,
    :language_lesson_member_id,
    :role,
    :body,
    :created_at
end
