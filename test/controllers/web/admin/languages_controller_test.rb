# frozen_string_literal: true

require 'test_helper'

class Web::Admin::LanguagesControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  test 'index' do
    get admin_languages_url
    assert_response :success
  end

  test 'new' do
    get new_admin_language_url
    assert_response :success
  end

  test 'create' do
    slug = 'racket'

    params = { language: { slug: slug } }
    post admin_languages_url, params: params
    assert_response :redirect

    assert { Language.find_by(slug: slug) }
  end

  test 'edit' do
    language = languages(:php)

    get edit_admin_language_url(language)
    assert_response :success
  end

  test 'update' do
    language = languages(:php)

    patch admin_language_url(language), params: { language: { progress: 'in_development' } }
    assert_response :redirect

    language.reload
    assert { language.progress.in_development? }
  end
end
