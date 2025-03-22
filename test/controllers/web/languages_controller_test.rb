require "test_helper"

class Web::LanguagesControllerTest < ActionDispatch::IntegrationTest
  test "show" do
    landing_page = language_landing_pages("javascript-ru")

    get language_url(landing_page.slug)
    assert_response :success
  end

  test "show (signed in)" do
    landing_page = language_landing_pages("javascript-ru")

    sign_in_as(:full)

    get language_url(landing_page.slug)
    assert_response :success
  end

  test "success (signed in)" do
    landing_page = language_landing_pages("elixir-ru")

    sign_in_as(:full)

    get language_url(landing_page.slug)
    assert_response :success
  end
end
