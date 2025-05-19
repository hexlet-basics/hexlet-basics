# == Schema Information
#
# Table name: surveys
#
#  id                                :bigint           not null, primary key
#  description                       :string
#  locale                            :string
#  question                          :string
#  run_after_finishing_lessons_count :integer          default(0)
#  run_always                        :boolean          default(FALSE)
#  slug                              :string
#  state                             :string
#  created_at                        :datetime         not null
#  updated_at                        :datetime         not null
#  parent_survey_id                  :bigint
#  parent_survey_item_id             :bigint
#
# Indexes
#
#  index_surveys_on_parent_survey_id       (parent_survey_id)
#  index_surveys_on_parent_survey_item_id  (parent_survey_item_id)
#  index_surveys_on_slug_and_locale        (slug,locale) UNIQUE
#
# Foreign Keys
#
#  fk_rails_...  (parent_survey_id => surveys.id)
#  fk_rails_...  (parent_survey_item_id => survey_items.id)
#
class Survey < ApplicationRecord
  has_many :items, class_name: "Survey::Item", dependent: :restrict_with_exception
  has_many :answers, class_name: "Survey::Answer", dependent: :restrict_with_exception
  has_many :scenario_items, class_name: "Survey::Scenario::Item"
  has_many :scenarios, through: :scenario_items, source: :scenario

  belongs_to :parent_survey_item, class_name: "Survey::Item", optional: true

  validates :question, presence: true
  validates :locale, presence: true
  validates :slug, presence: true, uniqueness: { scope: :locale }
  validates :items, length: { minimum: 2 }

  def self.ransackable_attributes(auth_object = nil)
    [ "created_at", "description", "id", "locale", "question", "slug", "state", "updated_at" ]
  end

  # def self.find_or_request_answer_if_needed_by(survey, user)
  #   # parent_survey_item = survey.parent_survey_item
  #   # # Rails.logger.info("LessonFinishedEvent #{parent_survey_item}")
  #   # if parent_survey_item
  #   #   # Rails.logger.info("LessonFinishedEvent PARENT")
  #   #   needed_answer = Survey::Answer.fulfilled.find_by user: user, survey_item: parent_survey_item
  #   #   return unless needed_answer
  #   # end
  #
  #   answer = survey.answers.find_or_initialize_by user: user
  #   answer.save! if answer.new_record?
  #   answer
  # end

  # def self.next_for_user(user, course_member: nil, ignore_filters: false)
  #   answers = user.survey_answers
  #   finished_lessons_count = course_member ? course_member.lesson_members.size : -1
  #
  #   # scope = where(parent_survey_item_id: answers.select(:survey_item_id))
  #   #   .where.not(id: answers.select(:survey_id))
  #   scope = where("parent_survey_item_id IS NULL OR parent_survey_item_id IN (?)", answers.select(:survey_item_id))
  #     .where.not(id: answers.select(:survey_id))
  #
  #   return scope if ignore_filters
  #
  #   scope = scope.where("run_always = ? OR run_after_finishing_lessons_count <= ?", true, finished_lessons_count)
  #   scope
  # end
end
