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
#  ahoy_visit_id       :bigint
#  user_id             :bigint           not null
#  ym_client_id        :string
#
# Indexes
#
#  index_leads_on_ahoy_visit_id  (ahoy_visit_id)
#  index_leads_on_user_id        (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (ahoy_visit_id => ahoy_visits.id)
#  fk_rails_...  (user_id => users.id)
#
require "test_helper"

class LeadTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
