require "test_helper"

class Web::Admin::MessagesControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = sign_in_as(:admin)
  end

  def test_index
    get admin_messages_url
    assert_response :success
  end
end
