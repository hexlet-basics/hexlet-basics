# frozen_string_literal: true

FactoryBot.define do
  factory :language_repository_version do
    language_name { 'MyString' }
    language { nil }
  end
end
