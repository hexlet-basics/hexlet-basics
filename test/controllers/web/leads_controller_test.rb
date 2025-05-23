require "test_helper"

class Web::LeadsControllerTest < ActionDispatch::IntegrationTest
  test "new" do
    sign_in_as(:full)

    get new_lead_url
    assert_response :success
  end

  test "create" do
    user = sign_in_as(:should_add_contact_method)

    assert { user.tag_list.include?("should_be_lead") }
    # lead_params = FactoryBot.attributes_for(:lead)
    lead_params = {
      contact_method: "telegram",
      contact_value: Faker::PhoneNumber.cell_phone
    }

    post leads_url, params: { lead: lead_params }
    assert_response :redirect

    user.reload
    # assert { user.survey_scenario_members.started.count == 1 }

    lead = user.leads.find_by! telegram: lead_params[:contact_value]
    # assert { user.survey_answers.requested.count == 2 }

    assert { user.tag_list.exclude?("should_be_lead") }

    assert { lead.courses_data.include?({ slug: "elixir", lessons_finished_count: 1 }) }
    assert { lead.survey_answers_data.include?({ question: "Какую задачу вы решаете?", answer: "Планирую поменять карьеру" }) }
  end
end
