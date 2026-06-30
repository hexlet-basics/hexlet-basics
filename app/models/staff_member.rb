# typed: strict
# frozen_string_literal: true

class StaffMember < ApplicationRecord
  belongs_to :user
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
