require "test_helper"

class Web::Admin::LanguageLessonMembersControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = sign_in_as(:admin)
  end

  def test_index
    get admin_language_lesson_members_path
    assert_response :success
  end
end
