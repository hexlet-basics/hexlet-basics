# frozen_string_literal: true

require 'test_helper'

class Api::Partners::YandexMarket::LanguagesControllerTest < ActionDispatch::IntegrationTest
  test '#index' do
    get api_partners_yandex_market_languages_url(format: :yml)
    assert_response :success
  end
end
