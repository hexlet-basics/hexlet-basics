# frozen_string_literal: true

# == Schema Information
#
# Table name: language_version_infos
#
#  id                  :integer          not null, primary key
#  description         :string
#  header              :string
#  keywords            :string
#  locale              :string
#  seo_description     :text
#  title               :string
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#  language_id         :bigint           not null
#  language_version_id :bigint           not null
#
# Indexes
#
#  index_language_version_infos_on_language_id          (language_id)
#  index_language_version_infos_on_language_version_id  (language_version_id)
#
# Foreign Keys
#
#  language_id          (language_id => languages.id)
#  language_version_id  (language_version_id => language_versions.id)
#
class Language::Version::Info < ApplicationRecord
  include Language::Version::InfoRepository

  serialize :keywords, type: Array

  validates :description, presence: true
  validates :header, presence: true

  belongs_to :language
  belongs_to :language_version, class_name: "Language::Version"

  def to_s
    header
  end
end
