# frozen_string_literal: true

# == Schema Information
#
# Table name: language_version_infos
#
#  id                  :bigint           not null, primary key
#  description         :string
#  header              :string
#  keywords            :string
#  locale              :string
#  seo_description     :text
#  title               :string
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#  language_id         :bigint           not null
#  language_version_id :bigint           not null
#
# Indexes
#
#  index_language_version_infos_on_language_id          (language_id)
#  index_language_version_infos_on_language_version_id  (language_version_id)
#
# Foreign Keys
#
#  fk_rails_...  (language_id => languages.id)
#  fk_rails_...  (language_version_id => language_versions.id)
#
require 'test_helper'

class Language::Version::InfoTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
