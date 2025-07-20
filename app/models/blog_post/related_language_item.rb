class BlogPost::RelatedLanguageItem < ApplicationRecord
  belongs_to :blog_post
  belongs_to :language

  validates :blog_post, uniqueness: { scope: :language }
end
