# frozen_string_literal: true

FactoryBot.define do
  factory :language_module_version do
    language_version { nil }
    order { 'MyString' }
  end
end
