# frozen_string_literal: true

FactoryBot.define do
  factory :review do
    language { nil }
    user { nil }
    state { 'MyString' }
    body { 'MyText' }
  end
end
