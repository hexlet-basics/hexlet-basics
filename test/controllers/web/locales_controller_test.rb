# frozen_string_literal: true

require 'test_helper'

class Web::LocalesControllerTest < ActionDispatch::IntegrationTest
  test 'switch to ru' do
    get switch_locale_url(subdomain: subdomain), params: {
      locale: 'ru'
    }
    # assert_response :redirect
    assert_redirected_to root_url(subdomain: :ru)
  end

  test 'switch to en' do
    get switch_locale_url(subdomain: subdomain), params: {
      locale: 'en'
    }
    # assert_response :redirect
    assert_redirected_to root_url(subdomain: nil)
  end

  test 'switch to unavailable locale' do
    get switch_locale_url(subdomain: subdomain), params: {
      locale: 'wrong'
    }
    assert_response :redirect
  end
end
