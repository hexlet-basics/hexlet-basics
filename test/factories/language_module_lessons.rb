# frozen_string_literal: true

FactoryBot.define do
  factory :language_module_lesson do
    order { 1 }
    natural_order { 1 }
    slug { "MyString" }
    original_code { "MyString" }
    prepared_code { "MyString" }
    test_code { "MyString" }
    path_to_code { "MyString" }
    language { nil }
    language_module { nil }
  end
end
