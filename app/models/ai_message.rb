# typed: strict
# frozen_string_literal: true

# == Schema Information
#
# Table name: ai_messages
#
#  id                    :bigint           not null, primary key
#  cache_creation_tokens :integer
#  cached_tokens         :integer
#  content               :text
#  content_raw           :json
#  input_tokens          :integer
#  output_tokens         :integer
#  role                  :string           not null
#  thinking_signature    :text
#  thinking_text         :text
#  thinking_tokens       :integer
#  created_at            :datetime         not null
#  updated_at            :datetime         not null
#  ai_chat_id            :bigint           not null
#  ai_model_id           :bigint
#  ai_tool_call_id       :bigint
#  user_id               :bigint
#
# Indexes
#
#  index_ai_messages_on_ai_chat_id       (ai_chat_id)
#  index_ai_messages_on_ai_model_id      (ai_model_id)
#  index_ai_messages_on_ai_tool_call_id  (ai_tool_call_id)
#  index_ai_messages_on_role             (role)
#  index_ai_messages_on_user_id          (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (ai_chat_id => ai_chats.id)
#  fk_rails_...  (ai_model_id => ai_models.id)
#  fk_rails_...  (ai_tool_call_id => ai_tool_calls.id)
#  fk_rails_...  (user_id => users.id)
#
class AiMessage < ApplicationRecord
  extend T::Sig

  acts_as_message chat: :ai_chat, tool_calls: :ai_tool_calls, model: :ai_model
  has_many_attached :attachments

  belongs_to :user, optional: true, counter_cache: :assistant_messages_count

  before_validation :assign_chat_user, on: :create

  scope :role_user, -> { where(role: "user") }

  sig { params(_auth_object: T.untyped).returns(T::Array[String]) }
  def self.ransackable_attributes(_auth_object = nil)
    [ "id", "role", "content", "user_id", "created_at", "ai_chat_id" ]
  end

  sig { params(_auth_object: T.untyped).returns(T::Array[String]) }
  def self.ransackable_associations(_auth_object = nil)
    [ "ai_chat" ]
  end

  private

  # RubyLLM persists messages without knowing about our user association,
  # so attribute the user-authored message to the chat owner ourselves.
  sig { void }
  def assign_chat_user
    return unless role == "user"
    return if user.present?

    self.user = ai_chat.user
  end
end
