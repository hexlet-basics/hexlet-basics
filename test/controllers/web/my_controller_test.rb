require "test_helper"

class Web::MyControllerTest < ActionDispatch::IntegrationTest
  def setup
    @user = users :full
    sign_in_as(:full)
  end

  test "should get index" do
    get my_url
    assert_response :success
  end
end
