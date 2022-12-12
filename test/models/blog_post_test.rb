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
require 'test_helper'

class BlogPostTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
