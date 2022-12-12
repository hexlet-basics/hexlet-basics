# frozen_string_literal: true

# == Schema Information
#
# Table name: language_modules
#
#  id          :bigint           not null, primary key
#  order       :integer
#  slug        :string(255)
#  state       :string(255)
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#  language_id :bigint
#  upload_id   :bigint
#
# Indexes
#
#  language_modules_language_id_index  (language_id)
#  language_modules_upload_id_index    (upload_id)
#
# Foreign Keys
#
#  language_modules_language_id_fkey  (language_id => languages.id)
#  language_modules_upload_id_fkey    (upload_id => uploads.id)
#
class Language::Module < ApplicationRecord
  belongs_to :language

  has_many :lessons, dependent: :destroy
  has_many :infos, through: :versions, class_name: 'Language::Module::Version::Info'
end
