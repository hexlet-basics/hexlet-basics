require "test_helper"

class Web::Admin::SurveysControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = sign_in_as(:admin)
  end

  def test_index
    get admin_surveys_url
    assert_response :success
  end

  def test_new
    get new_admin_survey_url
    assert_response :success
  end

  def test_create
    attrs = attributes_for(:survey, locale: I18n.locale)
    attrs[:items_attributes] = [
      attributes_for("survey/item"),
      attributes_for("survey/item")
    ]
    post admin_surveys_url, params: { survey: attrs }
    assert_response :redirect

    assert { Survey.find_by question: attrs[:question] }
  end

  def test_edit
    survey = surveys("goal")

    get edit_admin_survey_url(survey)
    assert_response :success
  end

  def test_update
    survey = surveys("goal")

    attrs = { question: "another-question" }
    # attrs[:items_attributes] = [
    #   attributes_for("survey/item"),
    #   attributes_for("survey/item")
    # ]

    patch admin_survey_url(survey), params: { survey: attrs }
    assert_response :redirect

    survey.reload
    assert { survey.question == attrs[:question] }
  end
end
