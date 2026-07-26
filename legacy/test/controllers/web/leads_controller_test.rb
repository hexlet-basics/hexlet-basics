# typed: true

require "test_helper"

class Web::LeadsControllerTest < ActionDispatch::IntegrationTest
  def test_new
    sign_in_as(:full)

    get new_lead_url

    assert_response :success
  end

  def test_create
    base_url = DepsLocator.current.amocrm.base_url.to_s

    WebMock.stub_request(:post, "#{base_url}/api/v4/leads/unsorted/forms")
      .to_return(
        status: 200,
        body: '{"_total_items":1}',
        headers: { "Content-Type" => "application/json" }
      )

    user = sign_in_as(:should_add_contact_method)

    assert { user.tag_list.include?("should_be_lead") }
    lead_params = {
      contact_method: "telegram",
      contact_value: Faker::PhoneNumber.cell_phone
    }

    post leads_url, params: { utm_source: "jopa", data: lead_params }

    assert_response :redirect

    user.reload

    lead = user.leads.find_by! telegram: lead_params[:contact_value]

    assert { user.tag_list.exclude?("should_be_lead") }

    assert { lead.courses_data.include?({ slug: "elixir", lessons_finished_count: 3 }) }
    assert { lead.survey_answers_data == [] }
  end
end
