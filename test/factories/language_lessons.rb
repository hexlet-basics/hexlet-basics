# frozen_string_literal: true

FactoryBot.define do
  factory :language_lesson do
    slug { 'MyString' }
    language { nil }
    language_module { nil }
    current_version { nil }
  end
end
