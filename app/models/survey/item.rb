# == Schema Information
#
# Table name: survey_items
#
#  id         :bigint           not null, primary key
#  order      :integer          not null
#  slug       :string
#  state      :string
#  value      :string
#  created_at :datetime         not null
#  updated_at :datetime         not null
#  survey_id  :integer          not null
#
# Indexes
#
#  index_survey_items_on_survey_id  (survey_id)
#
# Foreign Keys
#
#  fk_rails_...  (survey_id => surveys.id)
#
class Survey::Item < ApplicationRecord
  acts_as_taggable_on :tags

  belongs_to :survey

  validates :value, presence: true, if: :active?

  enum :state, { active: "active", archived: "archived" }, default: "active"
end
