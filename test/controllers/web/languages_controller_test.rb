# frozen_string_literal: true

require "test_helper"

class Web::LanguagesControllerTest < ActionDispatch::IntegrationTest
  test "show" do
    language = languages(:javascript)

    get language_url(language.slug)
    assert_response :success
  end

  test "show (signed in)" do
    language = languages(:javascript)

    sign_in_as(:full)

    get language_url(language.slug)
    assert_response :success
  end
end
