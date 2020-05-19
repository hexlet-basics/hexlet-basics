FactoryBot.define do
  factory :user do
    sequence(:email) { |n| "#{n}@email.com" }
    password { "password" }
  end
end
