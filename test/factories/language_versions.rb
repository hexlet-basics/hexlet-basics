# frozen_string_literal: true

# == Schema Information
#
# Table name: language_versions
#
#  id                     :bigint           not null, primary key
#  docker_image           :string
#  exercise_filename      :string
#  exercise_test_filename :string
#  extension              :string
#  name                   :string
#  state                  :string
#  result                 :string
#  language_id            :bigint           not null
#  created_at             :datetime         not null
#  updated_at             :datetime         not null
#  learn_as               :string
#  progress               :string
#
FactoryBot.define do
  factory :language_version do
    docker_image { 'MyString' }
    exercise_filename { 'MyString' }
    exercise_test_filename { 'MyString' }
    extension { 'MyString' }
    name { 'MyString' }
  end
end
