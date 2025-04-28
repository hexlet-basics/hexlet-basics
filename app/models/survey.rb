# == Schema Information
#
# Table name: surveys
#
#  id          :bigint           not null, primary key
#  description :string
#  locale      :string
#  question    :string
#  slug        :string
#  state       :string
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#
# Indexes
#
#  index_surveys_on_slug_and_locale  (slug,locale) UNIQUE
#
class Survey < ApplicationRecord
  has_many :items, class_name: "Survey::Item", dependent: :restrict_with_exception
  has_many :answers, class_name: "Survey::Answer", dependent: :restrict_with_exception

  validates :question, presence: true
  validates :locale, presence: true
  validates :slug, presence: true, uniqueness: { scope: :locale }
  validates :items, length: { minimum: 2 }

  def self.ransackable_attributes(auth_object = nil)
    [ "created_at", "description", "id", "locale", "question", "slug", "state", "updated_at" ]
  end

  def self.find_or_request_answer_by(slug, user)
    survey = self.find_by slug: slug
    return unless survey

    answer = survey.answers.find_or_initialize_by user: user
    answer.save! if answer.new_record?
    answer
  end
end
