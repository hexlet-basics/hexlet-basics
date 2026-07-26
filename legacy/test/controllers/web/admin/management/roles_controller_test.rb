# typed: true
# frozen_string_literal: true

require "test_helper"

class Web::Admin::Management::RolesControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
    @role = StaffMember::Role.create!(name: "Редактор")
  end

  def test_index
    get admin_management_roles_url

    assert_response :success
  end

  def test_new
    get new_admin_management_role_url

    assert_response :success
  end

  def test_create
    post admin_management_roles_url, params: { data: { name: "Маркетолог", description: "desc" } }

    assert_response :redirect
    assert { StaffMember::Role.exists?(name: "Маркетолог") }
  end

  def test_create_invalid
    post admin_management_roles_url, params: { data: { name: "" } }

    assert_response :redirect
  end

  def test_edit
    get edit_admin_management_role_url(@role)

    assert_response :success
  end

  def test_update
    patch admin_management_role_url(@role), params: { data: { name: "Главред" } }

    assert_response :redirect
    assert { @role.reload.name == "Главред" }
  end

  def test_destroy
    delete admin_management_role_url(@role)

    assert_response :redirect
    assert { !StaffMember::Role.exists?(@role.id) }
  end

  def test_permissions_show
    get admin_management_role_permission_url(@role)

    assert_response :success
  end

  def test_permissions_sync
    patch admin_management_role_permission_url(@role), params: {
      data: { permissions: [ { resource: "blog_posts", can_index: true, can_create: true } ] }
    }

    assert_response :redirect
    perm = @role.permissions.find_by(resource: "blog_posts")
    assert { perm.can_index? }
    assert { perm.can_create? }
    assert { !perm.can_destroy? }
  end
end
