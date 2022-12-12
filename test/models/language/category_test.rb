# frozen_string_literal: true

# == Schema Information
#
# Table name: language_categories
#
#  id         :bigint           not null, primary key
#  name_ru    :string
#  name_en    :string
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
