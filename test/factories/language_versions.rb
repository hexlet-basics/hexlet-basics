# frozen_string_literal: true

# == Schema Information
#
# Table name: language_versions
#
#  id                     :integer          not null, primary key
#  docker_image           :string
#  exercise_filename      :string
#  exercise_test_filename :string
#  extension              :string
#  learn_as               :string
#  name                   :string
#  progress               :string
#  result                 :string
#  state                  :string
#  created_at             :datetime         not null
#  updated_at             :datetime         not null
#  language_id            :bigint           not null
#
# Indexes
#
#  index_language_versions_on_language_id  (language_id)
#
# Foreign Keys
#
#  language_id  (language_id => languages.id)
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
