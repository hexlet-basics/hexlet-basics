# frozen_string_literal: true

require "test_helper"

class Web::Admin::BannersControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = sign_in_as(:admin)
  end

  def test_index
    get admin_banners_url
    assert_response :success
  end

  def test_new
    get new_admin_banner_url
    assert_response :success
  end

  def test_create
    attrs = attributes_for(:banner)
    post admin_banners_url, params: { data: attrs }
    assert_response :redirect

    assert { Banner.find_by(body: attrs[:body]) }
  end

  def test_edit
    banner = banners("published-ru")

    get edit_admin_banner_url(banner)
    assert_response :success
  end

  def test_update
    banner = banners("published-ru")

    assert { banner.published_state? }

    patch admin_banner_url(banner), params: { data: { body: "mumu", state: "archived" } }
    assert_response :redirect

    banner.reload
    assert { banner.archived_state? }
    assert { banner.body == "mumu" }
  end

  def test_destroy
    banner = banners("published-ru")

    delete admin_banner_url(banner)
    assert_response :redirect

    assert { Banner.find_by(id: banner.id).nil? }
  end
end
