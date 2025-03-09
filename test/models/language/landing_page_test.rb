# == Schema Information
#
# Table name: language_landing_pages
#
#  id                   :integer          not null, primary key
#  description          :string
#  header               :string
#  locale               :string
#  main                 :boolean
#  meta_description     :string
#  meta_title           :string
#  order                :string
#  slug                 :string
#  state                :string
#  created_at           :datetime         not null
#  updated_at           :datetime         not null
#  language_category_id :integer
#  language_id          :integer          not null
#
# Indexes
#
#  index_language_landing_pages_on_language_category_id  (language_category_id)
#  index_language_landing_pages_on_language_id           (language_id)
#
# Foreign Keys
#
#  language_category_id  (language_category_id => language_categories.id)
#  language_id           (language_id => languages.id)
#
require "test_helper"

class Language::LandingPageTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
