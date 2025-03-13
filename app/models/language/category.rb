# == Schema Information
#
# Table name: language_categories
#
#  id         :integer          not null, primary key
#  locale     :string
#  name       :string
#  name_en    :string
#  name_ru    :string
#  slug       :string
#  created_at :datetime         not null
#  updated_at :datetime         not null
#
class Language::Category < ApplicationRecord
  include Language::CategoryRepository

  validates :name, presence: true

  has_many :languages, ->(category) { where(category: category) }, dependent: :nullify, inverse_of: :category
  has_many :language_landing_pages, class_name: "Language::LandingPage", foreign_key: "language_category_id"
  has_many :language_versions, through: :languages, source: :versions
  has_many :language_version_infos, through: :language_versions, source: :infos, class_name: "Language::Version::Info"
  has_many :blog_posts, through: :languages, dependent: :restrict_with_exception

  def self.ransackable_attributes(auth_object = nil)
    []
  end

  def self.ransackable_associations(auth_object = nil)
    []
  end

  def name
    send :"name_#{I18n.locale}"
  end

  def to_s
    name
  end
end
