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
class BlogPost::Like < ApplicationRecord
  belongs_to :blog_post
  belongs_to :user, optional: true
end
