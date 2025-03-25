# frozen_string_literal: true

# == Schema Information
#
# Table name: language_lesson_members
#
#  id                 :integer          not null, primary key
#  state              :string
#  created_at         :datetime         not null
#  updated_at         :datetime         not null
#  language_id        :bigint           not null
#  language_member_id :integer          not null
#  lesson_id          :bigint           not null
#  openai_thread_id   :string
#  user_id            :bigint           not null
#
# Indexes
#
#  index_language_lesson_members_on_language_member_id            (language_member_id)
#  user_finished_lessons_language_module_lesson_id_index          (lesson_id)
#  user_finished_lessons_user_id_index                            (user_id)
#  user_finished_lessons_user_id_language_module_lesson_id_index  (user_id,lesson_id) UNIQUE
#
# Foreign Keys
#
#  language_member_id  (language_member_id => language_members.id)
#  lesson_id           (lesson_id => language_lessons.id)
#  user_id             (user_id => users.id)
#
class Language::Lesson::Member < ApplicationRecord
  include AASM
  include Language::Lesson::MemberRepository

  belongs_to :user
  belongs_to :language
  belongs_to :language_member, class_name: "Language::Member"
  belongs_to :lesson

  counter_culture :language_member,
                  column_name: ->(model) { model.finished? ? "finished_lessons_count" : nil }

  aasm :state do
    state :started, initial: true
    state :finished

    event :finish do
      transitions from: %i[started], to: :finished
    end
  end

  def serializable_data
    attributes.extract! "id", "state", "created_at"
  end
end
