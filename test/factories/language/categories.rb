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
FactoryBot.define do
  factory :language_category, class: 'Language::Category' do
    name { 'MyString' }
  end
end
