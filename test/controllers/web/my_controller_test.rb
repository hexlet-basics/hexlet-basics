require "test_helper"

class Web::MyControllerTest < ActionDispatch::IntegrationTest
  def setup
    @user = users :full
    sign_in_as(:full)
  end

  def test_should_get_index
    get my_url
    assert_response :success
  end
end
