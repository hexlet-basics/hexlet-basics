# frozen_string_literal: true

FactoryBot.define do
  factory :language_module_description do
    name { 'MyString' }
    description { 'MyText' }
    locale { 'MyString' }
    language_module { nil }
    language { '' }
  end
end
