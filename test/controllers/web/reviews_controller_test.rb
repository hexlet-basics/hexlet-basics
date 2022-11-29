# frozen_string_literal: true

require 'test_helper'

class ReviewsControllerTest < ActionDispatch::IntegrationTest
  test 'index' do
    get reviews_url
    assert_response :success
  end

  test 'show' do
    review = reviews('full-javascript')
    get review_url(review)
    assert_response :success
  end
end
