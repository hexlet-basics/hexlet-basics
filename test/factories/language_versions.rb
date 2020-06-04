# frozen_string_literal: true

FactoryBot.define do
  factory :language_version do
    docker_image { 'MyString' }
    exercise_filename { 'MyString' }
    exercise_test_filename { 'MyString' }
    extension { 'MyString' }
    name { 'MyString' }
  end
end
