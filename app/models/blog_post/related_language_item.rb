# == Schema Information
#
# Table name: blog_post_related_language_items
#
#  id           :bigint           not null, primary key
#  order        :integer
#  created_at   :datetime         not null
#  updated_at   :datetime         not null
#  blog_post_id :bigint           not null
#  language_id  :bigint           not null
#
# Indexes
#
#  index_blog_post_related_language_items_on_blog_post_id  (blog_post_id)
#  index_blog_post_related_language_items_on_language_id   (language_id)
#
# Foreign Keys
#
#  fk_rails_...  (blog_post_id => blog_posts.id)
#  fk_rails_...  (language_id => languages.id)
#
class BlogPost::RelatedLanguageItem < ApplicationRecord
  belongs_to :blog_post
  belongs_to :language

  validates :blog_post, uniqueness: { scope: :language }
end
