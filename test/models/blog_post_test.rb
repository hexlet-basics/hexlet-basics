# frozen_string_literal: true

# == Schema Information
#
# Table name: blog_posts
#
#  id          :integer          not null, primary key
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
#  creator_id   (creator_id => users.id)
#  language_id  (language_id => languages.id)
#
require "test_helper"

class BlogPostTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
