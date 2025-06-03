# frozen_string_literal: true

require "test_helper"

class Web::CasesControllerTest < ActionDispatch::IntegrationTest
  test "index" do
    get cases_url
    assert_response :success
  end

  test "for teachers cases" do
    get for_teachers_cases_url
    assert_response :success
  end
end
