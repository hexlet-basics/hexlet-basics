require "test_helper"

class Web::Admin::SurveyAnswersControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = sign_in_as(:admin)
  end

  test "index" do
    get admin_survey_answers_url
    assert_response :success
  end
end
