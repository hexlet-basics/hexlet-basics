require "test_helper"

class Web::LanguagesControllerTest < ActionDispatch::IntegrationTest
  def test_show
    landing_page = language_landing_pages("javascript-ru")

    get language_url(landing_page.slug)
    assert_response :success
  end

  def test_show_archived
    archived_landing_page = language_landing_pages("javascript-ru-archived")
    landing_page = language_landing_pages("javascript-ru")

    get language_url(archived_landing_page.slug)
    assert_redirected_to language_path(landing_page.slug)
  end

  def test_show_signed_in
    landing_page = language_landing_pages("javascript-ru")

    sign_in_as(:full)

    get language_url(landing_page.slug)
    assert_response :success
  end

  def test_show_non_member_signed_in
    landing_page = language_landing_pages("php-ru")

    sign_in_as(:full)

    get language_url(landing_page.slug)
    assert_response :success
  end

  def test_success_signed_in
    landing_page = language_landing_pages("elixir-ru")

    sign_in_as(:full)

    get language_url(landing_page.slug)
    assert_response :success
  end
end
