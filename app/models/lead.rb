# == Schema Information
#
# Table name: leads
#
#  id                  :bigint           not null, primary key
#  courses_data        :text
#  email               :string
#  phone               :string
#  state               :string
#  survey_answers_data :text
#  telegram            :string
#  whatsapp            :string
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#  user_id             :bigint           not null
#
# Indexes
#
#  index_leads_on_user_id  (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (user_id => users.id)
#
class Lead < ApplicationRecord
  belongs_to :user
  serialize :survey_answers_data
  serialize :courses_data

  def self.ransackable_attributes(auth_object = nil)
    [ "created_at", "email", "id", "phone", "state", "telegram", "updated_at", "user_id", "whatsapp" ]
  end
end
