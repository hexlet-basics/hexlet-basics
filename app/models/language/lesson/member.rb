# frozen_string_literal: true

class Language::Lesson::Member < ApplicationRecord
  include AASM
  include Language::Lesson::MemberRepository

  counter_culture :language

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
