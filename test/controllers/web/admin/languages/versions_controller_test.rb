# frozen_string_literal: true

require 'test_helper'

class Web::Admin::Languages::VersionsControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  test 'index' do
    language = languages(:php)

    get admin_language_versions_url(language, subdomain: subdomain)

    assert_response :success
  end

  test 'create' do
    language = languages(:php)

    post admin_language_versions_url(language, subdomain: subdomain)
    assert_response :redirect
    language.reload

    language_module = language.modules.find_by(slug: 'basics')
    module_version = language.current_module_versions.find_by(module: language_module)
    lesson = language.lessons.find_by(slug: 'hello-world')
    lesson_version = language.current_lesson_versions.find_by(lesson: lesson)

    assert { language_module }
    assert { module_version }
    assert { module_version.infos.find_by(locale: :ru) }
    assert { module_version.infos.find_by(locale: :en) }
    assert { lesson }
    assert { lesson_version }
    assert { lesson_version.infos.find_by(locale: :ru) }
    assert { lesson_version.infos.find_by(locale: :en) }
  end
end
