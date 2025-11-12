require "test_helper"

class Web::BooksControllerTest < ActionDispatch::IntegrationTest
  def test_show
    get book_path
    assert_response :success
  end

  def test_show_sign_in
    sign_in_as(:full)

    get book_path
    assert_response :success
  end

  def test_create_request
    user = sign_in_as(:ready_to_start_learning)
    assert { user.book_request.nil? }

    post create_request_book_path
    assert_response :redirect

    user.reload

    assert { user.book_request }
    assert { user.survey_scenario_members.count == 3 }
  end

  def test_download
    user = sign_in_as(:full)

    get download_book_path
    assert_response :success

    assert { user.book_request.downloaded_state? }
  end
end
