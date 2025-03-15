require "test_helper"

class Web::MyControllerTest < ActionDispatch::IntegrationTest
  test "should get index" do
    get my_url
    assert_response :success
  end
end
