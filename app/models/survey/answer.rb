# == Schema Information
#
# Table name: survey_answers
#
#  id             :bigint           not null, primary key
#  state          :string
#  created_at     :datetime         not null
#  updated_at     :datetime         not null
#  survey_id      :integer          not null
#  survey_item_id :integer
#  user_id        :integer          not null
#
# Indexes
#
#  index_survey_answers_on_survey_id              (survey_id)
#  index_survey_answers_on_survey_id_and_user_id  (survey_id,user_id) UNIQUE
#  index_survey_answers_on_survey_item_id         (survey_item_id)
#  index_survey_answers_on_user_id                (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (survey_id => surveys.id)
#  fk_rails_...  (survey_item_id => survey_items.id)
#  fk_rails_...  (user_id => users.id)
#
class Survey::Answer < ApplicationRecord
  include AASM

  belongs_to :survey
  belongs_to :survey_item, class_name: "Survey::Item", optional: true
  belongs_to :user

  validates :survey, uniqueness: { scope: :user }

  # enum :state, { requested: "requested", fulfilled: "fulfilled" }, default: "requested", suffix: true

  def self.ransackable_attributes(auth_object = nil)
    [ "created_at", "id", "state", "survey_id", "survey_item_id", "updated_at", "user_id" ]
  end

  # aasm :state, enum: true do
  #   state :requested
  #   state :fulfilled do
  #     validates :survey_item, optional: false
  #   end
  # end
end
