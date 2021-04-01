# frozen_string_literal: true

FactoryBot.define do
  factory :user do
    sequence(:email) { |n| "Email-#{n}@email.com" }
    password { 'password' }
  end
end
