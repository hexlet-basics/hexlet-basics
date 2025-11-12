# frozen_string_literal: true

require "test_helper"

class Web::CasesControllerTest < ActionDispatch::IntegrationTest
  def test_index
    get cases_url
    assert_response :success
  end

  def test_for_teachers_cases
    get for_teachers_cases_url
    assert_response :success
  end
end
