# frozen_string_literal: true

require 'test_helper'

class Web::LocalesControllerTest < ActionDispatch::IntegrationTest
  test 'switch' do
    get switch_locale_url(subdomain: subdomain), params: {
      locale: 'ru'
    }
    assert_response :redirect
  end

  test 'switch to unavailable locale' do
    get switch_locale_url(subdomain: subdomain), params: {
      locale: 'wrong'
    }
    assert_response :redirect
  end
end
