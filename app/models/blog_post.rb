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
class BlogPost < ApplicationRecord
  include BlogPostRepository

  has_one_attached :cover do |attachable|
    attachable.variant :main, resize_to_limit: [ 1020, 480 ], preprocessed: true
    attachable.variant :list, resize_to_limit: [ 510, 240 ], preprocessed: true
    attachable.variant :thumb, resize_to_limit: [ 205, 120 ], preprocessed: true
  end

  def self.ransackable_attributes(_auth_object = nil)
    [ "created_at" ]
  end

  def self.ransackable_associations(_auth_object = nil)
    []
  end

  enum :state, { draft: "draft", published: "published", archived: "archived" }, suffix: true, validate: true

  validates :name, presence: true
  validates :slug, presence: true, uniqueness: true, format: { with: /\A[\w-]+\z/ }
  validates :locale, presence: true
  validates :body, presence: true
  validates :description, presence: true

  belongs_to :language, optional: true
  has_one :category, through: :language, class_name: "Language::Category"
  belongs_to :creator, class_name: "User"

  def to_s
    name
  end
end
