# frozen_string_literal: true

FactoryBot.define do
  factory :blog_post do
    course { nil }
    locale { 'MyString' }
    state { 'MyString' }
    slug { 'MyString' }
    name { 'MyString' }
    body { 'MyText' }
  end
end
