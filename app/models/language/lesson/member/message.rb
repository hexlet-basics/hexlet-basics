# == Schema Information
#
# Table name: language_lesson_member_messages
#
#  id                        :integer          not null, primary key
#  body                      :text
#  role                      :string
#  created_at                :datetime         not null
#  updated_at                :datetime         not null
#  language_id               :integer          not null
#  language_lesson_id        :integer          not null
#  language_lesson_member_id :integer          not null
#
# Indexes
#
#  idx_on_language_lesson_member_id_fe254654e9                  (language_lesson_member_id)
#  index_language_lesson_member_messages_on_language_id         (language_id)
#  index_language_lesson_member_messages_on_language_lesson_id  (language_lesson_id)
#
# Foreign Keys
#
#  language_id                (language_id => languages.id)
#  language_lesson_id         (language_lesson_id => language_lessons.id)
#  language_lesson_member_id  (language_lesson_member_id => language_lesson_members.id)
#
class Language::Lesson::Member::Message < ApplicationRecord
  belongs_to :language, class_name: "Language"
  belongs_to :language_lesson, class_name: "Language::Lesson"
  belongs_to :language_lesson_member, class_name: "Language::Lesson::Member"

  validates :role, presence: true
  validates :body, presence: true

  enum :role, { user: "user", assistant: "assistant" }, suffix: true, validate: true

  def self.ransackable_attributes(auth_object = nil)
    []
  end

  def self.ransackable_associations(auth_object = nil)
    []
  end
end
