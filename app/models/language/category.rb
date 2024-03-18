# frozen_string_literal: true

# == Schema Information
#
# Table name: language_categories
#
#  id         :bigint           not null, primary key
#  name_ru    :string
#  name_en    :string
#  slug       :string
#  created_at :datetime         not null
#  updated_at :datetime         not null
#
class Language::Category < ApplicationRecord
  validates :name, presence: true

  has_many :languages, ->(category) { where(category: category) }, dependent: :nullify, inverse_of: :category
  has_many :blog_posts, through: :languages, dependent: :restrict_with_exception

  def name
    respond_to?("name_#{I18n.locale}") ? send("name_#{I18n.locale}") : name_en
  end

  def to_s
    name
  end
end
