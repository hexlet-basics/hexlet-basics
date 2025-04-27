# == Schema Information
#
# Table name: surveys
#
#  id          :integer          not null, primary key
#  description :string
#  locale      :string
#  question    :string
#  slug        :string
#  state       :string
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#
# Indexes
#
#  index_surveys_on_slug_and_locale  (slug,locale) UNIQUE
#
FactoryBot.define do
  factory :survey do
    question { Faker::Lorem.question }
    slug { Faker::Internet.slug }
    # locale { :ru }
    description { Faker::Lorem.paragraph_by_chars }
  end
end
