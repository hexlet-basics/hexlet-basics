# frozen_string_literal: true

# == Schema Information
#
# Table name: language_module_version_infos
#
#  id                  :bigint           not null, primary key
#  name                :string
#  description         :string
#  locale              :string
#  language_id         :bigint           not null
#  version_id          :bigint           not null
#  language_version_id :bigint           not null
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#
class Language::Module::Version::Info < ApplicationRecord
  include Language::Module::Version::InfoRepository

  belongs_to :language
  belongs_to :language_version, class_name: 'Language::Version'
  belongs_to :version

  def to_s
    name
  end
end
