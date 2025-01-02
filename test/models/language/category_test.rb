# frozen_string_literal: true

# == Schema Information
#
# Table name: language_categories
#
#  id         :integer          not null, primary key
#  name_en    :string
#  name_ru    :string
#  slug       :string
#  created_at :datetime         not null
#  updated_at :datetime         not null
#
require 'test_helper'

class Language::CategoryTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
