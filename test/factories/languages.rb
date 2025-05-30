# frozen_string_literal: true

# == Schema Information
#
# Table name: languages
#
#  id                          :bigint           not null, primary key
#  docker_image                :string(255)
#  exercise_filename           :string(255)
#  exercise_test_filename      :string(255)
#  extension                   :string(255)
#  hexlet_program_landing_page :string
#  learn_as                    :string
#  lessons_count               :integer          default(0), not null
#  members_count               :integer          default(0), not null
#  name                        :string(255)
#  order                       :integer
#  progress                    :string
#  slug                        :string(255)
#  state                       :string(255)
#  created_at                  :datetime         not null
#  updated_at                  :datetime         not null
#  category_id                 :bigint
#  current_version_id          :bigint
#  openai_assistant_id         :string
#  upload_id                   :bigint
#
# Indexes
#
#  index_languages_on_category_id         (category_id)
#  index_languages_on_current_version_id  (current_version_id)
#  languages_slug_index                   (slug) UNIQUE
#  languages_upload_id_index              (upload_id)
#
# Foreign Keys
#
#  fk_rails_...  (category_id => language_categories.id)
#  fk_rails_...  (current_version_id => language_versions.id)
#  fk_rails_...  (upload_id => uploads.id)
#
FactoryBot.define do
  factory :language do
    slug { "php" }
  end
end
