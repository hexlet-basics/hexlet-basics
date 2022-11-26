require 'test_helper'

class Web::Admin::ReviewsControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  test 'index' do
    get admin_reviews_url
    assert_response :success
  end
end
