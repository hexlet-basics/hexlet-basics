require "test_helper"

class Api::FeedsControllerTest < ActionDispatch::IntegrationTest
  test "should get index" do
    get yandex_courses_api_feeds_path(format: "xml")
    assert_response :success
  end
end
