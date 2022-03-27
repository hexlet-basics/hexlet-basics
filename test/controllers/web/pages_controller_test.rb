# frozen_string_literal: true

require 'test_helper'

class Web::PagesControllerTest < ActionDispatch::IntegrationTest
  %w[tos about privacy].each do |slug|
    test slug do
      get page_url(slug, subdomain: subdomain)
      assert_response :success
    end
  end
end
