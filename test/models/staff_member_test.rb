# typed: true
# frozen_string_literal: true

require "test_helper"

# == Schema Information
#
# Table name: staff_members
#
#  id              :bigint           not null, primary key
#  allowed_locales :string           default(["ru"]), not null, is an Array
#  created_at      :datetime         not null
#  updated_at      :datetime         not null
#  role_id         :bigint           not null
#  user_id         :bigint           not null
#
# Indexes
#
#  index_staff_members_on_role_id  (role_id)
#  index_staff_members_on_user_id  (user_id) UNIQUE
#
# Foreign Keys
#
#  fk_rails_...  (role_id => staff_member_roles.id)
#  fk_rails_...  (user_id => users.id)
#
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
