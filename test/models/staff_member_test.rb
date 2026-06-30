# typed: true
# frozen_string_literal: true

require "test_helper"

class StaffMemberTest < ActiveSupport::TestCase
  setup do
    @role = StaffMember::Role.create!(name: "Редактор")
    @role.permissions.create!(resource: "blog_posts", can_index: true, can_create: true)
    @staff_member = StaffMember.create!(user: users(:one), role: @role)
  end

  def test_resource_permission_predicates
    assert { @staff_member.can_index_resource?("blog_posts") }
    assert { @staff_member.can_create_resource?("blog_posts") }
    assert { !@staff_member.can_update_resource?("blog_posts") }
    assert { !@staff_member.can_destroy_resource?("blog_posts") }
    assert { !@staff_member.can_index_resource?("banners") }
  end

  def test_can_any_in
    assert { @staff_member.can_any_in?([ StaffMember::Role::Permission::Resource::BlogPosts ]) }
    assert { !@staff_member.can_any_in?([ StaffMember::Role::Permission::Resource::Banners ]) }
  end

  def test_user_uniqueness
    duplicate = StaffMember.new(user: users(:one), role: @role)
    assert { !duplicate.valid? }
  end

  def test_user_staff_predicate
    assert { users(:one).reload.staff? }
    assert { !users(:two).reload.staff? }
    assert { users(:admin).staff? }
  end

  def test_permission_uniqueness_per_role
    duplicate = @role.permissions.build(resource: "blog_posts")
    assert { !duplicate.valid? }
  end
end
