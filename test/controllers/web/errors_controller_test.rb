require "test_helper"

class Web::ErrorsControllerTest < ActionDispatch::IntegrationTest
  def test_404
    get "/unknownpage"
    assert_response :not_found
  end
end
