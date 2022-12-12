# frozen_string_literal: true

# == Schema Information
#
# Table name: language_lesson_versions
#
#  id                  :bigint           not null, primary key
#  order               :integer
#  original_code       :string
#  prepared_code       :string
#  test_code           :string
#  natural_order       :integer
#  path_to_code        :string
#  language_version_id :bigint           not null
#  language_id         :bigint           not null
#  lesson_id           :bigint           not null
#  module_version_id   :bigint           not null
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#
FactoryBot.define do
  factory :language_lesson_version do
    # TODO: implement later
  end
end
