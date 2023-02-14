# frozen_string_literal: true

require 'test_helper'

class Web::PagesControllerTest < ActionDispatch::IntegrationTest
  %w[tos about privacy authors].each do |slug|
    test slug do
      get page_url(slug)
      assert_response :success
    end
  end
end
