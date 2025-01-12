# frozen_string_literal: true

require "test_helper"

class Web::LocalesControllerTest < ActionDispatch::IntegrationTest
  test "switch to ru" do
    get switch_locale_url(locale: nil), params: {
      new_locale: :ru
    }
    assert_redirected_to root_url(suffix: :ru)
  end

  test "switch to en" do
    get switch_locale_url(locale: :ru), params: {
      new_locale: :en
    }
    assert_redirected_to root_url(suffix: nil)
  end

  test "switch to unavailable locale" do
    get switch_locale_url, params: {
      new_locale: "wrong"
    }
    assert_response :redirect
  end
end
