require "test_helper"

class Web::Scenarios::SurveysControllerTest < ActionDispatch::IntegrationTest
  test "show new scenario" do
    user = sign_in_as(:just_signed_up)

    base_scenario = survey_scenarios("base")
    member = user.survey_scenario_members.find_by scenario: base_scenario

    get scenario_survey_path(base_scenario, member.next_survey)
    assert_response :success
  end

  test "show inside scenario" do
    user = sign_in_as(:inside_base_survey_scenario)

    base_scenario = survey_scenarios("base")
    member = user.survey_scenario_members.find_by scenario: base_scenario

    get scenario_survey_path(base_scenario, member.next_survey)
    assert_response :success
  end

  test "show without member" do
    user = sign_in_as(:admin)

    members = user.survey_scenario_members
    assert { members.count == 0 }

    base_scenario = survey_scenarios("base")

    get scenario_survey_path(base_scenario, base_scenario.surveys.first)
    assert_response :redirect
  end

  test "show if answered" do
    sign_in_as(:full)

    base_scenario = survey_scenarios("base")

    get scenario_survey_path(base_scenario, base_scenario.surveys.first)
    assert_response :redirect
  end
end
