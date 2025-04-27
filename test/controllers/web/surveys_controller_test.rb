require "test_helper"

class Web::SurveysControllerTest < ActionDispatch::IntegrationTest
  test "#show" do
    user = sign_in_as(:just_signed_up)

    survey = surveys("goal")

    get survey_path(survey.slug)
    assert_response :success
  end
end
