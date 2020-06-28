# frozen_string_literal: true

FactoryBot.define do
  factory :lesson_member do
    user { nil }
    language { nil }
    state { 'MyString' }
    language_version { nil }
  end
end
