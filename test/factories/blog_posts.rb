# frozen_string_literal: true

# == Schema Information
#
# Table name: blog_posts
#
#  id          :bigint           not null, primary key
#  language_id :bigint
#  locale      :string
#  state       :string
#  slug        :string
#  name        :string
#  body        :text
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#  creator_id  :bigint           not null
#  description :string
#
FactoryBot.define do
  factory :blog_post do
    locale { I18n.locale }
    slug { Faker::Internet.slug }
    name { Faker::Name.first_name }
    body { Faker::Lorem.paragraph_by_chars }
    description { Faker::Lorem.paragraph }
  end
end
