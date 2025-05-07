require "test_helper"

class Web::Scenarios::Surveys::AnswersControllerTest < ActionDispatch::IntegrationTest
  test "create" do
    user = sign_in_as(:just_signed_up)

    base_scenario = survey_scenarios("base")
    member = user.survey_scenario_members.find_by scenario: base_scenario
    assert { user.survey_scenario_members.count == 1 }
    survey = member.next_survey
    assert { !user.survey_answers.find_by(survey: survey) }

    attrs = {
      survey_answer: {
        survey_item_id: survey.items.first.id
      }
    }
    post scenario_survey_answers_path(base_scenario, survey), params: attrs
    assert_response :redirect

    assert { user.survey_answers.find_by(survey: survey) }
    assert { survey != member.next_survey }
  end

  test "create last survey" do
    user = sign_in_as(:inside_base_survey_scenario)

    base_scenario = survey_scenarios("base")
    member = user.survey_scenario_members.find_by scenario: base_scenario
    survey = member.next_survey

    item = survey.items.find_by! slug: "goal-item1"
    attrs = {
      survey_answer: {
        survey_item_id: item.id
      }
    }
    post scenario_survey_answers_path(base_scenario, survey), params: attrs
    assert_response :redirect

    member.reload
    user.reload

    assert { member.state == "finished" }
    assert { !member.next_survey }
    assert { user.survey_scenario_members.count == 1 }
  end
end
