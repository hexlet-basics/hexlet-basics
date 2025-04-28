# frozen_string_literal: true

# == Schema Information
#
# Table name: reviews
#
#  id          :bigint           not null, primary key
#  body        :text
#  first_name  :string
#  last_name   :string
#  locale      :string
#  pinned      :boolean
#  state       :string
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#  language_id :bigint           not null
#  user_id     :bigint           not null
#
# Indexes
#
#  index_reviews_on_language_id  (language_id)
#  index_reviews_on_user_id      (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (language_id => languages.id)
#  fk_rails_...  (user_id => users.id)
#
FactoryBot.define do
  factory :review do
    first_name { Faker::Name.first_name }
    last_name { Faker::Name.last_name }
    state { "published" }
    body { Faker::Lorem.paragraph_by_chars }
  end
end
