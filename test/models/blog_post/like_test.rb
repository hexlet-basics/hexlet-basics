# == Schema Information
#
# Table name: blog_post_likes
#
#  id           :bigint           not null, primary key
#  created_at   :datetime         not null
#  updated_at   :datetime         not null
#  blog_post_id :bigint           not null
#  user_id      :bigint
#
# Indexes
#
#  index_blog_post_likes_on_blog_post_id  (blog_post_id)
#  index_blog_post_likes_on_user_id       (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (blog_post_id => blog_posts.id)
#  fk_rails_...  (user_id => users.id)
#
require "test_helper"

class BlogPost::LikeTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
