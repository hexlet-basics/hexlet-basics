# frozen_string_literal: true

# == Schema Information
#
# Table name: language_lesson_members
#
#  id                 :bigint           not null, primary key
#  user_id            :bigint           not null
#  lesson_id          :bigint           not null
#  updated_at         :datetime         not null
#  state              :string
#  language_id        :bigint           not null
#  created_at         :datetime         not null
#  language_member_id :integer          not null
#
class Language::Lesson::Member < ApplicationRecord
  include AASM
  include Language::Lesson::MemberRepository

  belongs_to :user
  belongs_to :language
  belongs_to :language_member, class_name: 'Language::Member'
  belongs_to :lesson

  counter_culture :language_member,
                  column_name: ->(model) { model.finished? ? 'finished_lessons_count' : nil }

  aasm :state do
    state :started, initial: true
    state :finished

    event :finish do
      transitions from: %i[started finished], to: :finished
    end
  end

  def to_hash(*_args)
    attrs = attributes.extract! 'id', 'state', 'created_at'
    attrs.to_hash
  end
end
