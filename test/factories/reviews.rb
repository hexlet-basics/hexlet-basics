# frozen_string_literal: true

# == Schema Information
#
# Table name: reviews
#
#  id          :bigint           not null, primary key
#  language_id :bigint           not null
#  user_id     :bigint           not null
#  state       :string
#  body        :text
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#  locale      :string
#
FactoryBot.define do
  factory :review do
    language { nil }
    user { nil }
    state { 'MyString' }
    body { 'MyText' }
  end
end
