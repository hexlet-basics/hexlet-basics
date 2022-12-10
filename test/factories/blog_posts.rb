# frozen_string_literal: true

FactoryBot.define do
  factory :blog_post do
    locale { I18n.locale }
    slug { Faker::Internet.slug }
    name { Faker::Name.first_name }
    body { Faker::Lorem.paragraph_by_chars }
    description { Faker::Lorem.paragraph }
  end
end
