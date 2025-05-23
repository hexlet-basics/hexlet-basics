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
  include AASM

  acts_as_taggable_on :tags

  belongs_to :survey

  enum :state, { active: "active", archived: "archived" }

  aasm :state, enum: true do
    state :active, initial: true do
      validates :value, presence: true
    end
    state :archived
  end
end
