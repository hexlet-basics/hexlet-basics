# frozen_string_literal: true

# == Schema Information
#
# Table name: language_modules
#
#  id          :bigint           not null, primary key
#  slug        :string(255)
#  state       :string(255)
#  order       :integer
#  language_id :bigint
#  upload_id   :bigint
#  updated_at  :datetime         not null
#  created_at  :datetime         not null
#
class Language::Module < ApplicationRecord
  belongs_to :language

  has_many :lessons, dependent: :destroy
  has_many :infos, through: :versions, class_name: 'Language::Module::Version::Info'
end
