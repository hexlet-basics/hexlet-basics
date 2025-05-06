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
    user = sign_in_as(:ready_to_start_learning)
    assert { user.book_request.nil? }

    post create_request_book_path
    assert_response :redirect

    user.reload

    assert { user.book_request }
    assert { user.survey_answers.requested.count == 7 }
  end

  test "download" do
    user = sign_in_as(:full)

    get download_book_path
    assert_response :success

    assert { user.book_request.downloaded? }
  end
end
