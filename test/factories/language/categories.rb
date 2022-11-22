# frozen_string_literal: true

FactoryBot.define do
  factory :language_category, class: 'Language::Category' do
    name { 'MyString' }
  end
end
