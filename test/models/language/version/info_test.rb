# frozen_string_literal: true

# == Schema Information
#
# Table name: language_version_infos
#
#  id                  :bigint           not null, primary key
#  language_id         :bigint           not null
#  language_version_id :bigint           not null
#  locale              :string
#  description         :string
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#  title               :string
#  seo_description     :text
#  header              :string
#  keywords            :string
#
require 'test_helper'

class Language::Version::InfoTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
