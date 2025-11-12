require "test_helper"

class Api::FeedsControllerTest < ActionDispatch::IntegrationTest
  def test_should_get_index
    get yandex_courses_api_feeds_path(format: "xml")
    assert_response :success
  end
end
