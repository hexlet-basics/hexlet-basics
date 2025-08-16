# frozen_string_literal: true

# == Schema Information
#
# Table name: blog_posts
#
#  id                           :bigint           not null, primary key
#  body                         :text
#  description                  :string
#  locale                       :string
#  name                         :string
#  related_language_items_count :integer          default(0), not null
#  slug                         :string
#  state                        :string
#  created_at                   :datetime         not null
#  updated_at                   :datetime         not null
#  creator_id                   :bigint           not null
#  language_id                  :bigint
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

  has_many :related_language_items, dependent: :delete_all
  has_many :related_languages, through: :related_language_items, source: :language
  has_many :related_main_language_landing_pages, through: :related_languages, source: :main_landing_page

  has_one_attached :cover do |attachable| # 800 x 400
    attachable.variant :main, resize_to_limit: [ 688, 344 ], preprocessed: true
    attachable.variant :list, resize_to_limit: [ 546, 273 ], preprocessed: true
    attachable.variant :thumb, resize_to_limit: [ 200, 100 ], preprocessed: true
  end

  def self.ransackable_attributes(_auth_object = nil)
    [ "id", "created_at", "related_language_items_count", "state" ]
  end

  def self.ransackable_associations(_auth_object = nil)
    []
  end

  enum :state, { draft: "draft", published: "published", archived: "archived" }, suffix: true, validate: true, default: "draft"

  validates :name, presence: true
  validates :slug, presence: true, uniqueness: true, format: { with: /\A[\w-]+\z/ }
  validates :locale, presence: true
  validates :body, presence: true
  validates :description, presence: true

  belongs_to :language, optional: true
  has_one :category, through: :language, class_name: "Language::Category"
  belongs_to :creator, class_name: "User"
  has_many :likes

  def to_s
    name
  end
end
