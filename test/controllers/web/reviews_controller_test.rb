# frozen_string_literal: true

require "test_helper"

class ReviewsControllerTest < ActionDispatch::IntegrationTest
  def test_index
    get reviews_url
    assert_response :success
  end
end
