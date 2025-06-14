require "test_helper"

class Web::LanguagesControllerTest < ActionDispatch::IntegrationTest
  test "show" do
    landing_page = language_landing_pages("javascript-ru")

    get language_url(landing_page.slug)
    assert_response :success
  end

  test "show (archived)" do
    archived_landing_page = language_landing_pages("javascript-ru-archived")
    landing_page = language_landing_pages("javascript-ru")

    get language_url(archived_landing_page.slug)
    assert_redirected_to language_path(landing_page.slug)
  end

  test "show (signed in)" do
    landing_page = language_landing_pages("javascript-ru")

    sign_in_as(:full)

    get language_url(landing_page.slug)
    assert_response :success
  end

  test "show (non-member/signed in)" do
    landing_page = language_landing_pages("php-ru")

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
