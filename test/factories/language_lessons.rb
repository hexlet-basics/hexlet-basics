# frozen_string_literal: true

# == Schema Information
#
# Table name: language_lessons
#
#  id            :bigint           not null, primary key
#  slug          :string(255)
#  state         :string(255)
#  order         :integer
#  original_code :text
#  prepared_code :text
#  test_code     :text
#  path_to_code  :string(255)
#  module_id     :bigint
#  language_id   :bigint
#  upload_id     :bigint
#  updated_at    :datetime         not null
#  natural_order :integer
#  created_at    :datetime         not null
#
FactoryBot.define do
  factory :language_lesson do
    slug { 'MyString' }
    language { nil }
    language_module { nil }
    current_version { nil }
  end
end
