# typed: true
# frozen_string_literal: true

require "test_helper"

class Web::Admin::StaffAccessTest < ActionDispatch::IntegrationTest
  setup do
    @editor = users(:one)
    role = StaffMember::Role.create!(name: "Редактор")
    role.permissions.create!(resource: "blog_posts", can_index: true, can_create: true, can_update: true)
    StaffMember.create!(user: @editor, role:)
  end

  def test_editor_can_access_blog
    sign_in_as(@editor)
    get admin_blog_posts_url
    assert_response :success
  end

  def test_editor_cannot_access_other_resource
    sign_in_as(@editor)
    get admin_banners_url
    assert_redirected_to admin_root_path
  end

  def test_editor_cannot_access_management
    sign_in_as(@editor)
    get admin_management_users_url
    assert_redirected_to admin_root_path
  end

  # Кастомный мутирующий экшен на чужом ресурсе должен быть запрещён (else false).
  def test_editor_cannot_post_custom_action_on_other_resource
    sign_in_as(@editor)
    post review_admin_language_url(languages(:php))
    assert_redirected_to admin_root_path
  end

  def test_non_staff_user_cannot_enter_admin
    sign_in_as(users(:two))
    get admin_blog_posts_url
    assert_redirected_to root_path
  end

  def test_admin_can_access_everything
    sign_in_as(:admin)
    get admin_blog_posts_url
    assert_response :success
    get admin_banners_url
    assert_response :success
    get admin_management_users_url
    assert_response :success
  end
end
