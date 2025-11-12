require "test_helper"

class Web::Admin::LanguageLandingPagesControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  def test_index
    get admin_language_landing_pages_url
    assert_response :success
  end

  def test_new
    get new_admin_language_landing_page_url
    assert_response :success
  end

  def test_create
    slug = "racket"

    params = { language_landing_page: { slug: slug } }
    post admin_language_landing_pages_url, params: params
    assert_response :redirect

    # assert { Language.find_by(slug: slug) }
  end

  def test_edit
    landing_page = language_landing_pages("php-ru")

    get edit_admin_language_landing_page_url(landing_page)
    assert_response :success
  end

  def test_update
    landing_page = language_landing_pages("php-ru")

    params = { language_landing_page: { description: "new description" } }
    patch admin_language_landing_page_url(landing_page), params: params
    assert_response :redirect

    # language.reload
    # assert { landing_page.progress.in_development? }
  end
end
