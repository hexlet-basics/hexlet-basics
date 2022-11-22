# frozen_string_literal: true

require 'test_helper'

class ReviewsControllerTest < ActionDispatch::IntegrationTest
  test 'should get index' do
    get reviews_url
    assert_response :success
  end

  test 'should get show' do
    get reviews_url
    assert_response :success
  end
end
