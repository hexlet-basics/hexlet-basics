# == Schema Information
#
# Table name: survey_scenarios
#
#  id             :bigint           not null, primary key
#  locale         :string
#  name           :string
#  state          :string
#  created_at     :datetime         not null
#  updated_at     :datetime         not null
#  survey_item_id :bigint
#
# Indexes
#
#  index_survey_scenarios_on_survey_item_id  (survey_item_id)
#
# Foreign Keys
#
#  fk_rails_...  (survey_item_id => survey_items.id)
#
class Survey::Scenario < ApplicationRecord
  belongs_to :survey_item, class_name: "Survey::Item", optional: true
  has_many :items, class_name: "Survey::Scenario::Item"
  has_many :surveys, through: :items
  has_many :members, class_name: "Survey::Scenario::Member"
  has_many :user, through: :members
  has_many :triggers

  validates :locale, presence: true
  validates :name, presence: true

  enum :state, { active: "active", archived: "archived" }, suffix: true

  def self.ransackable_attributes(auth_object = nil)
    [ "created_at", "id", "locale", "name", "state", "survey_item_id", "updated_at" ]
  end
end
