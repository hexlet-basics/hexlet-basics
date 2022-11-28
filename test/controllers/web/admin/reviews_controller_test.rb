# frozen_string_literal: true

require 'test_helper'

class Web::Admin::ReviewsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = sign_in_as(:admin)
  end

  test 'index' do
    get admin_reviews_url
    assert_response :success
  end

  test 'new' do
    get new_admin_review_url
    assert_response :success
  end

  test 'create' do
    lang = languages(:php)

    params = { review: {
      user_id: @user.id, language_id: lang.id, body: 'ehu', locale: I18n.locale
    } }
    post admin_reviews_url, params: params
    assert_response :redirect

    assert { lang.reviews.find_by body: params[:review][:body] }
  end

  test 'edit' do
    review = reviews('full-javascript')

    get edit_admin_review_url(review)
    assert_response :success
  end

  test 'update' do
    review = reviews('full-javascript')

    patch admin_review_url(review), params: { review: { body: 'mumu' } }
    assert_response :redirect

    review.reload
    assert { review.body == 'mumu' }
  end
end
