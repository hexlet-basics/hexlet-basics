# frozen_string_literal: true

require 'test_helper'

class Web::Admin::LanguagesControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  test 'index' do
    get admin_languages_path
    assert_response :success
  end

  test 'new' do
    get new_admin_language_path
    assert_response :success
  end

  test 'create' do
    slug = 'php'

    post admin_languages_path, params: { language: { slug: slug } }
    assert_response :redirect

    assert { Language.find_by(slug: slug) }
  end

  test 'edit' do
    language = languages(:one)

    get edit_admin_language_path(language)
    assert_response :success
  end

  test 'update' do
    language = languages(:one)

    patch admin_language_path(language), params: { language: { progress: 'in_development' } }
    assert_response :redirect

    language.reload
    assert { language.progress.in_development? }
  end
end
