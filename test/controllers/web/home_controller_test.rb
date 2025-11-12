require "test_helper"

class Web::HomeControllerTest < ActionDispatch::IntegrationTest
  def test_index
    get root_url
    assert_response :success
  end

  def test_index_sign_in
    sign_in_as(:full)

    get root_url
    assert_response :success
  end

  def test_index_with_stored_locale
    open_session do |s|
      s.get s.root_url(suffix: :ru)
      s.assert_response :success

      s.get s.root_url(suffix: nil)
      s.assert_response :redirect
    end
  end

  def test_map
    get map_url
    assert_response :success
  end
end
