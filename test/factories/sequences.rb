# frozen_string_literal: true

FactoryBot.define do
  sequence :github_auth_hash do
    {
      provider: 'github',
      uid: '550e8400-e29b-41d4-a716-446655440000',
      info: {
        email: Faker::Internet.email.upcase,
        first_name: Faker::Name.first_name,
        last_name: Faker::Name.last_name,
        nickname: Faker::Internet.username
      }
    }
  end
end
