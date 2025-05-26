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
#  ym_client_id        :string
#
# Indexes
#
#  index_leads_on_user_id  (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (user_id => users.id)
#
FactoryBot.define do
  factory :lead do
    user { nil }
    email { Faker::Internet.email }
    phone { Faker::PhoneNumber.cell_phone }
    telegram { Faker::Internet.username }
    whatsapp { Faker::PhoneNumber.cell_phone }
    survey_answers_data { [] }
  end
end
