require "test_helper"

class Web::Admin::LanguageLandingPagesControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  test "index" do
    get admin_language_landing_pages_url
    assert_response :success
  end

  test "new" do
    get new_admin_language_landing_page_url
    assert_response :success
  end

  test "create" do
    slug = "racket"

    params = { language_landing_page: { slug: slug } }
    post admin_language_landing_pages_url, params: params
    assert_response :redirect

    # assert { Language.find_by(slug: slug) }
  end

  test "edit" do
    landing_page = language_landing_pages("php-ru")

    get edit_admin_language_landing_page_url(landing_page)
    assert_response :success
  end

  test "update" do
    landing_page = language_landing_pages("php-ru")

    params = { language_landing_page: { description: "new description" } }
    patch admin_language_landing_page_url(landing_page), params: params
    assert_response :redirect

    # language.reload
    # assert { landing_page.progress.in_development? }
  end
end
