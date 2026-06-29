# typed: false
# frozen_string_literal: true

# == Schema Information
#
# Table name: ai_chats
#
#  id                        :bigint           not null, primary key
#  created_at                :datetime         not null
#  updated_at                :datetime         not null
#  ai_model_id               :bigint
#  language_lesson_member_id :bigint           not null
#  user_id                   :bigint           not null
#
# Indexes
#
#  index_ai_chats_on_ai_model_id                            (ai_model_id)
#  index_ai_chats_on_language_lesson_member_id              (language_lesson_member_id)
#  index_ai_chats_on_user_id                                (user_id)
#  index_ai_chats_on_user_id_and_language_lesson_member_id  (user_id,language_lesson_member_id) UNIQUE
#
# Foreign Keys
#
#  fk_rails_...  (ai_model_id => ai_models.id)
#  fk_rails_...  (language_lesson_member_id => language_lesson_members.id)
#  fk_rails_...  (user_id => users.id)
#
class AiChat < ApplicationRecord
  acts_as_chat messages: :ai_messages, model: :ai_model

  belongs_to :user
  belongs_to :language_lesson_member, class_name: "Language::Lesson::Member"

  validates :user, uniqueness: { scope: :language_lesson_member }
end
