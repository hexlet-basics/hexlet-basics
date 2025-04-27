require "test_helper"

class Web::Surveys::AnswersControllerTest < ActionDispatch::IntegrationTest
  test "#create" do
    user = sign_in_as(:just_signed_up)

    survey = surveys("goal")
    answer = user.survey_answers.find_by survey: survey
    assert { answer.requested? }

    attrs = {
      survey_answer: {
        survey_item_id: survey.items.first.id
      }
    }
    post survey_answers_path(survey.slug), params: attrs
    assert_response :redirect

    answer.reload
    assert { answer.fulfilled? }
  end
end
