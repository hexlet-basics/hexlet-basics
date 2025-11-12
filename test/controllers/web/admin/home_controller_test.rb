require "test_helper"

class Web::Admin::HomeControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  def test_index
    get admin_root_url
    assert_response :success
  end
end
