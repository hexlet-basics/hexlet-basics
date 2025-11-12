require "test_helper"

class Web::Admin::SurveyScenariosControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = sign_in_as(:admin)
  end

  def test_index
    get admin_survey_scenarios_url
    assert_response :success
  end

  def test_new
    get new_admin_survey_scenario_url
    assert_response :success
  end

  # test "create" do
  #   attrs = attributes_for(:survey, locale: I18n.locale)
  #   attrs[:items_attributes] = [
  #     attributes_for("survey/item"),
  #     attributes_for("survey/item")
  #   ]
  #   post admin_surveys_url, params: { survey: attrs }
  #   assert_response :redirect
  #
  #   assert { Survey.find_by question: attrs[:question] }
  # end
  #
  # test "edit" do
  #   survey = surveys("goal")
  #
  #   get edit_admin_survey_url(survey)
  #   assert_response :success
  # end
  #
  # test "update" do
  #   survey = surveys("goal")
  #
  #   attrs = { question: "another-question" }
  #   # attrs[:items_attributes] = [
  #   #   attributes_for("survey/item"),
  #   #   attributes_for("survey/item")
  #   # ]
  #
  #   patch admin_survey_url(survey), params: { survey: attrs }
  #   assert_response :redirect
  #
  #   survey.reload
  #   assert { survey.question == attrs[:question] }
  # end
end
