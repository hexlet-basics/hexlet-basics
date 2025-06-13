# frozen_string_literal: true

require "test_helper"

class Web::Admin::ReviewsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = sign_in_as(:admin)
  end

  test "index" do
    get admin_reviews_url
    assert_response :success
  end

  test "new" do
    get new_admin_review_url
    assert_response :success
  end

  test "create" do
    lang = languages(:php)

    attrs = attributes_for(:review, language_id: lang.id, user_id: @user.id, locale: I18n.locale)
    post admin_reviews_url, params: { review: attrs }
    assert_response :redirect

    assert { lang.reviews.find_by body: attrs[:body] }
  end

  test "edit" do
    review = reviews("full-javascript")

    get edit_admin_review_url(review)
    assert_response :success
  end

  test "update" do
    review = reviews("full-javascript")

    assert { review.published_state? }

    patch admin_review_url(review), params: { review: { body: "mumu", state_event: "archive" } }
    assert_response :redirect

    review.reload
    assert { review.archived? }
    assert { review.body == "mumu" }
  end
end
