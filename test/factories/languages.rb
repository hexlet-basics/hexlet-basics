# frozen_string_literal: true

# == Schema Information
#
# Table name: languages
#
#  id                     :bigint           not null, primary key
#  name                   :string(255)
#  slug                   :string(255)
#  extension              :string(255)
#  docker_image           :string(255)
#  exercise_filename      :string(255)
#  exercise_test_filename :string(255)
#  state                  :string(255)
#  upload_id              :bigint
#  updated_at             :datetime         not null
#  progress               :string
#  current_version_id     :bigint
#  created_at             :datetime         not null
#  learn_as               :string
#  lessons_count          :integer          default(0), not null
#  members_count          :integer          default(0), not null
#  order                  :integer
#  category_id            :bigint
#
FactoryBot.define do
  factory :language do
    slug { 'php' }
  end
end
