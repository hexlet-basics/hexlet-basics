require "test_helper"

class Web::ErrorsControllerTest < ActionDispatch::IntegrationTest
  test "404" do
    get "/unknownpage"
    assert_response :not_found
  end
end
