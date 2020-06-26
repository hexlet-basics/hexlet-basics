# frozen_string_literal: true

require 'test_helper'

class Web::Admin::Languages::UploadsControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  test 'create' do
    language = languages(:one)

    post admin_language_uploads_path(language.id)
    assert_response :redirect

    lesson_module = language.modules.find_by(slug: 'basics')
    lesson = language.lessons.find_by(slug: 'hello-world')

    assert { Language::Upload.find_by(language: language.id) }
    assert { lesson_module }
    assert { lesson_module.current_version }
    assert { lesson_module.descriptions.find_by(locale: :ru) }
    assert { lesson_module.descriptions.find_by(locale: :en) }
    assert { lesson }
    assert { lesson.current_version }
    assert { lesson.descriptions.find_by(locale: :ru) }
    assert { lesson.descriptions.find_by(locale: :en) }
  end
end
