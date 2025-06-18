class Language::Lesson::MemberResource < ApplicationResource
  typelize_from Language::Lesson::Member

  attributes :id, :user_id, :state, :openai_thread_id, :messages_count, :created_at

  typelize state: [ enum: [ "started", "finished" ] ]

  typelize :string
  attribute :language_lesson_slug do |obj|
    obj.lesson.slug
  end


   typelize :string
   attribute :language_slug do |obj|
     obj.language.slug
   end

  typelize :string
  attribute :language_lesson_name do |obj|
    info = obj.lesson.localed_info
    info&.name
  end
end
