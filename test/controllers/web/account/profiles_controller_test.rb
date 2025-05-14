# frozen_string_literal: true

require "test_helper"

class Web::Account::ProfilesControllerTest < ActionDispatch::IntegrationTest
  test "edit" do
    sign_in_as(:full)

    get edit_account_profile_url
    assert_response :success
  end

  test "update" do
    user = sign_in_as(:full)
    new_name = "new first name"

    patch account_profile_url(id: user.id), params: {
      user: {
        first_name: new_name,
        contact_method: :telegram,
        contact_value: "@orgprog"
      }
    }
    assert_response :redirect

    user.reload

    assert { user.first_name == new_name }
    assert { user.lead }
    assert { user.lead.courses_data.include?({ slug: "elixir", lessons_finished_count: 3 }) }
    assert { user.lead.survey_answers_data.include?({ question: "Какую задачу вы решаете?", answer: "Планирую поменять карьеру" }) }
  end

  test "update (new lead)" do
    user = sign_in_as("should_add_contact_method")
    patch account_profile_url(id: user.id), params: {
      user: {
        contact_method: :telegram,
        contact_value: "@orgprog"
      }
    }
    assert_response :redirect

    user.reload

    assert { user.lead }
    # assert { @user.lead.courses_data.include?({ slug: "elixir", lessons_finished_count: 3 }) }
    # assert { @user.lead.survey_answers_data.include?({ question: "Какую задачу вы решаете?", answer: "Планирую поменять карьеру" }) }
  end

  test "destroy" do
    user = sign_in_as(:full)
    delete account_profile_url
    assert_response :redirect

    user.reload

    assert { user.removed? }
    assert { !signed_in? }
  end
end
