# typed: true

require "test_helper"

class Web::UsersControllerTest < ActionDispatch::IntegrationTest
  extend T::Sig

  class SetupContext < T::Struct
    const :language, Language
    const :lesson, Language::Lesson
    const :version, Language::Lesson::Version
  end

  sig { returns(SetupContext) }
  def context
    language = languages(:ruby)
    lesson = language.lessons.first!
    SetupContext.new(language:, lesson:, version: lesson.versions.first!)
  end

  def test_new
    get new_user_url
    assert_response :success
  end

  def test_create
    user_params = FactoryBot.attributes_for(:user)

    post users_url, params: { data: user_params }
    assert_response :redirect

    user = User.find_by! email: user_params[:email].downcase

    assert { authenticated? }
    assert { user.survey_scenario_members.empty? }
  end

  def test_create_with_demo
    ctx = context
    user_params = FactoryBot.attributes_for(:user)

    code = file_fixture("exercise/correct.rb").read
    post check_api_lesson_url(ctx.lesson), params: { version_id: ctx.version.id, data: { attributes: { code: code } } }

    post users_url, params: { data: user_params }
    user = User.find_by! email: user_params[:email].downcase

    assert { ctx.language.members.exists? user: user }
    assert { ctx.lesson.members.finished.exists? user: user }
    assert { authenticated? }
  end
end
