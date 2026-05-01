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
#  survey_id  :bigint           not null
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
  class State < T::Enum
    enums do
      Active = new("active")
      Archived = new("archived")
    end
  end

  acts_as_taggable_on :tags

  belongs_to :survey

  validates :value, presence: true, if: :active?

  typed_enum :state, State, default: State::Active
end
