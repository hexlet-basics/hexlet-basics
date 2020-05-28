# frozen_string_literal: true

FactoryBot.define do
  factory :language_module_lesson_version do
    original_code { 'MyString' }
    prepared_code { 'MyString' }
    test_code { 'MyString' }
    path_to_code { 'MyString' }
    lesson { nil }
  end
end
