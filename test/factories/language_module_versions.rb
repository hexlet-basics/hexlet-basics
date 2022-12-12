# frozen_string_literal: true

# == Schema Information
#
# Table name: language_module_versions
#
#  id                  :bigint           not null, primary key
#  language_id         :bigint           not null
#  language_version_id :bigint           not null
#  module_id           :bigint           not null
#  order               :integer
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#
FactoryBot.define do
  factory :language_module_version do
    language_version { nil }
    order { 'MyString' }
  end
end
