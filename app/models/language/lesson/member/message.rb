# == Schema Information
#
# Table name: language_lesson_member_messages
#
#  id                        :bigint           not null, primary key
#  body                      :text
#  role                      :string
#  created_at                :datetime         not null
#  updated_at                :datetime         not null
#  language_id               :bigint           not null
#  language_lesson_id        :bigint           not null
#  language_lesson_member_id :bigint           not null
#  user_id                   :bigint
#
# Indexes
#
#  idx_on_language_lesson_member_id_fe254654e9                  (language_lesson_member_id)
#  index_language_lesson_member_messages_on_language_id         (language_id)
#  index_language_lesson_member_messages_on_language_lesson_id  (language_lesson_id)
#  index_language_lesson_member_messages_on_user_id             (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (language_id => languages.id)
#  fk_rails_...  (language_lesson_id => language_lessons.id)
#  fk_rails_...  (language_lesson_member_id => language_lesson_members.id)
#  fk_rails_...  (user_id => users.id)
#
class Language::Lesson::Member::Message < ApplicationRecord
  # TODO: add locale
  belongs_to :language, class_name: "Language"
  belongs_to :language_lesson, class_name: "Language::Lesson"
  belongs_to :language_lesson_member, class_name: "Language::Lesson::Member"
  belongs_to :user

  validates :role, presence: true
  # validates :body, presence: true

  counter_culture :language_lesson_member, column_name: "messages_count"
  counter_culture :user, column_name: "assistant_messages_count"

  enum :role, { user: "user", assistant: "assistant" }, suffix: true, validate: true

  def self.ransackable_attributes(auth_object = nil)
    [ "id", "created_at", "language_lesson_member_id" ]
  end

  def self.ransackable_associations(auth_object = nil)
    []
  end

  def to_s
    body
  end
end
