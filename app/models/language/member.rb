# frozen_string_literal: true

# == Schema Information
#
# Table name: language_members
#
#  id                     :integer          not null, primary key
#  finished_lessons_count :integer          default(0), not null
#  state                  :string
#  created_at             :datetime         not null
#  updated_at             :datetime         not null
#  language_id            :bigint           not null
#  user_id                :bigint           not null
#
# Indexes
#
#  index_language_members_on_language_id  (language_id)
#  index_language_members_on_user_id      (user_id)
#
# Foreign Keys
#
#  language_id  (language_id => languages.id)
#  user_id      (user_id => users.id)
#
class Language::Member < ApplicationRecord
  include AASM

  def self.ransackable_attributes(_auth_object = nil)
    %w[created_at language_id]
  end

  belongs_to :user
  belongs_to :language
  has_many :lesson_members, class_name: 'Language::Lesson::Member', foreign_key: :language_member_id, inverse_of: :language_member, dependent: :destroy

  counter_culture :language

  aasm :state do
    state :started, initial: true
    state :finished

    event :finish do
      transitions from: %i[started finished], to: :finished, guard: :all_lessons_finished?
    end
  end

  def serializable_data
    attributes.extract! 'id', 'state', 'created_at'
  end

  def all_lessons_finished?
    not_finished_lessons = user.not_finished_lessons_for_language(language)

    not_finished_lessons.empty?
  end
end
