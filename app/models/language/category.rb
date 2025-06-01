# == Schema Information
#
# Table name: language_categories
#
#  id          :bigint           not null, primary key
#  description :string
#  header      :string
#  locale      :string
#  name        :string
#  name_en     :string
#  name_ru     :string
#  slug        :string
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#
class Language::Category < ApplicationRecord
  include Language::CategoryRepository

  validates :name, presence: true, uniqueness: true
  validates :header, presence: true, uniqueness: true
  validates :slug, presence: true, uniqueness: { scope: :locale }

  has_many :items, class_name: "Language::Category::Item", foreign_key: "language_category_id"
  has_many :language_landing_pages, through: :items, source: "language_landing_page"
  # has_many :qna_items, foreign_key: "language_landing_page_id"

  has_many :languages, ->(category) { where(category: category) }, dependent: :nullify, inverse_of: :category
  has_many :language_versions, through: :languages, source: :versions
  has_many :language_version_infos, through: :language_versions, source: :infos, class_name: "Language::Version::Info"
  has_many :blog_posts, through: :languages, dependent: :restrict_with_exception

  def self.ransackable_attributes(auth_object = nil)
    []
  end

  def self.ransackable_associations(auth_object = nil)
    []
  end

  def to_s
    name
  end
end
