# typed: strict
# frozen_string_literal: true

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
class StaffMember < ApplicationRecord
  belongs_to :user, inverse_of: :staff_member
  belongs_to :role, class_name: "StaffMember::Role"

  validates :user, uniqueness: true

  sig { params(resource: String).returns(T::Boolean) }
  def can_index_resource?(resource)
    role.permissions.any? { |p| p.resource == resource && p.can_index? }
  end

  sig { params(resource: String).returns(T::Boolean) }
  def can_create_resource?(resource)
    role.permissions.any? { |p| p.resource == resource && p.can_create? }
  end

  sig { params(resource: String).returns(T::Boolean) }
  def can_update_resource?(resource)
    role.permissions.any? { |p| p.resource == resource && p.can_update? }
  end

  sig { params(resource: String).returns(T::Boolean) }
  def can_destroy_resource?(resource)
    role.permissions.any? { |p| p.resource == resource && p.can_destroy? }
  end

  sig { params(resources: T::Array[StaffMember::Role::Permission::Resource]).returns(T::Boolean) }
  def can_any_in?(resources)
    serialized = resources.map(&:serialize)
    role.permissions.any? { |p| serialized.include?(p.resource) }
  end

  sig { params(_auth_object: T.untyped).returns(T::Array[String]) }
  def self.ransackable_attributes(_auth_object = nil)
    %w[id user_id role_id]
  end

  sig { params(_auth_object: T.untyped).returns(T::Array[String]) }
  def self.ransackable_associations(_auth_object = nil)
    %w[user role]
  end
end
