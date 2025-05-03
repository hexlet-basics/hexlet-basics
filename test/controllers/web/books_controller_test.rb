require "test_helper"

class Web::BooksControllerTest < ActionDispatch::IntegrationTest
  test "show" do
    get book_path
    assert_response :success
  end

  test "show (sign in)" do
    sign_in_as(:full)

    get book_path
    assert_response :success
  end

  test "create_request" do
    _user = users(:ready_to_start_learning)

    post create_request_book_path
    assert_response :redirect
  end

  test "download" do
    sign_in_as(:full)

    get download_book_path
    assert_response :success
  end
end
