# typed: true
# frozen_string_literal: true

require "test_helper"

class Web::LocalesControllerTest < ActionDispatch::IntegrationTest
  def test_switch_to_ru
    get switch_locale_url(locale: nil), params: {
      new_locale: :ru
    }
    assert_redirected_to root_url(suffix: :ru)
  end

  def test_switch_to_en
    get switch_locale_url(locale: :ru), params: {
      new_locale: :en
    }
    assert_redirected_to root_url(suffix: nil)
  end

  def test_switch_to_unavailable_locale
    get switch_locale_url, params: {
      new_locale: "wrong"
    }
    assert_response :redirect
  end

  def test_switch_preserves_current_page
    get switch_locale_url, params: { new_locale: :en },
      headers: { "HTTP_REFERER" => languages_url(suffix: :ru) }
    assert_redirected_to languages_path(suffix: nil)
  end

  def test_switch_preserves_current_page_to_ru
    get switch_locale_url, params: { new_locale: :ru },
      headers: { "HTTP_REFERER" => languages_url(suffix: nil) }
    assert_redirected_to languages_path(suffix: :ru)
  end

  def test_switch_ignores_foreign_referer
    get switch_locale_url, params: { new_locale: :ru },
      headers: { "HTTP_REFERER" => "https://evil.example/ru/languages" }
    assert_redirected_to root_url(suffix: :ru)
  end
end
