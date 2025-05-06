require "test_helper"

class Web::UsersControllerTest < ActionDispatch::IntegrationTest
  test "new" do
    get new_user_url
  end

  test "create" do
    user_params = FactoryBot.attributes_for(:user)

    post users_url, params: { user_sign_up_form: user_params }
    assert_response :redirect

    user = User.find_by email: user_params[:email].downcase

    assert { user.present? }
    assert { signed_in? }
    assert { user.survey_answers.requested.count == 2 }
  end

  test "create (with demo)" do
    user_params = FactoryBot.attributes_for(:user)

    language = languages(:ruby)
    lesson = language.lessons.first
    code = file_fixture("exercise/correct.rb").read
    post check_api_lesson_url(lesson), params: { version_id: lesson.versions.first.id, data: { attributes: { code: code } } }

    post users_url, params: { user_sign_up_form: user_params }
    user = User.find_by email: user_params[:email].downcase

    assert { language.members.exists? user: user }
    assert { lesson.members.finished.exists? user: user }
    assert { user.present? }
    assert { signed_in? }
  end
end
