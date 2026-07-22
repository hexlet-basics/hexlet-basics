# typed: true
# frozen_string_literal: true

require "test_helper"

class Web::Admin::Management::StaffMembersControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
    @role = StaffMember::Role.create!(name: "Редактор")
    @user = users(:one)
  end

  def test_index
    get admin_management_staff_members_url

    assert_response :success
  end

  def test_new
    get new_admin_management_staff_member_url

    assert_response :success
  end

  def test_create
    post admin_management_staff_members_url, params: {
      data: { user_id: @user.id, role_id: @role.id, allowed_locales: [ "ru" ] }
    }

    assert_response :redirect
    assert { StaffMember.exists?(user: @user, role: @role) }
  end

  def test_update
    staff_member = StaffMember.create!(user: @user, role: @role)
    other_role = StaffMember::Role.create!(name: "Маркетолог")

    patch admin_management_staff_member_url(staff_member), params: {
      data: { role_id: other_role.id, allowed_locales: [ "ru", "en" ] }
    }

    assert_response :redirect
    assert { staff_member.reload.role == other_role }
  end

  def test_destroy
    staff_member = StaffMember.create!(user: @user, role: @role)
    delete admin_management_staff_member_url(staff_member)

    assert_response :redirect
    assert { !StaffMember.exists?(staff_member.id) }
  end
end
