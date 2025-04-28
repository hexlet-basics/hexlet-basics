# frozen_string_literal: true

# == Schema Information
#
# Table name: blog_posts
#
#  id          :bigint           not null, primary key
#  body        :text
#  description :string
#  locale      :string
#  name        :string
#  slug        :string
#  state       :string
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#  creator_id  :bigint           not null
#  language_id :bigint
#
# Indexes
#
#  index_blog_posts_on_creator_id   (creator_id)
#  index_blog_posts_on_language_id  (language_id)
#  index_blog_posts_on_slug         (slug) UNIQUE
#
# Foreign Keys
#
#  fk_rails_...  (creator_id => users.id)
#  fk_rails_...  (language_id => languages.id)
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
