require "test_helper"

class Web::UsersControllerTest < ActionDispatch::IntegrationTest
  def test_new
    get new_user_url
    assert_response :success
  end

  def test_create
    user_params = FactoryBot.attributes_for(:user)

    post users_url, params: { user: user_params }
    assert_response :redirect

    user = User.find_by email: user_params[:email].downcase

    assert { user.present? }
    assert { signed_in? }
    assert { user.survey_scenario_members.started.count == 1 }
    # assert { user.survey_answers.requested.count == 2 }
  end

  def test_create_with_demo
    user_params = FactoryBot.attributes_for(:user)

    language = languages(:ruby)
    lesson = language.lessons.first
    code = file_fixture("exercise/correct.rb").read
    post check_api_lesson_url(lesson), params: { version_id: lesson.versions.first.id, data: { attributes: { code: code } } }

    post users_url, params: { user: user_params }
    user = User.find_by email: user_params[:email].downcase

    assert { language.members.exists? user: user }
    assert { lesson.members.finished.exists? user: user }
    assert { user.present? }
    assert { signed_in? }
  end
end
