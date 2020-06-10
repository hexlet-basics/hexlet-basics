# frozen_string_literal: true

require 'test_helper'

class Web::Languages::Modules::LessonsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @language = Language.first
    @module = @language.modules.first
    @lesson = @module.lessons.first
  end

  test 'show' do
    get language_module_lesson_path(@language.slug, @module.slug, @lesson.slug)
    assert_response :success
  end
end
