# typed: strict

class AiMessageResource < ApplicationResource
  typelize_from AiMessage

  attributes :id, :role, :user_id, :created_at

  typelize :string
  attribute :content do
    it.content
  end

  typelize :string
  attribute :language_slug do
    it.ai_chat.language_lesson_member.lesson.language.slug
  end

  typelize :string
  attribute :language_lesson_slug do
    it.ai_chat.language_lesson_member.lesson.slug
  end

  typelize :string
  attribute :language_lesson_name do
    it.ai_chat.language_lesson_member.lesson.localed_info&.name
  end
end
