# frozen_string_literal: true

require 'test_helper'

class Web::Admin::Languages::VersionsControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  test 'index' do
    language = languages(:one)

    get admin_language_versions_path(language)

    assert_response :success
  end

  test 'create' do
    language = languages(:one)

    post admin_language_versions_path(language)
    assert_response :redirect

    lesson_module = language.modules.find_by(slug: 'basics')
    lesson = language.lessons.find_by(slug: 'hello-world')

    assert { Language::Version.find_by(language: language.id) }
    assert { lesson_module }
    assert { lesson_module.version }
    assert { lesson_module.version.data.find_by(locale: :ru) }
    assert { lesson_module.version.data.find_by(locale: :en) }
    assert { lesson }
    assert { lesson.version }
    assert { lesson.version.data.find_by(locale: :ru) }
    assert { lesson.version.data.find_by(locale: :en) }
  end
end
