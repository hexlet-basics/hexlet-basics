# frozen_string_literal: true

require 'test_helper'

class Web::HomeControllerTest < ActionDispatch::IntegrationTest

  test 'index' do
    get root_url
    assert_response :success
  end

  test '#index with stored locale' do
    open_session do |s|
      s.get s.root_url(locale: :ru), headers: {
        'User-Agent': 'Mozilla'
      }
      s.assert_response :success

      s.get s.root_url(locale: nil), headers: {
        'User-Agent': 'Mozilla'
      }
      s.assert_response :redirect
    end
  end

  test '#robots' do
    get robots_url(format: :txt)
    assert_response :success
  end
end
