# typed: strict
# frozen_string_literal: true

# == Schema Information
#
# Table name: staff_member_roles
#
#  id          :bigint           not null, primary key
#  description :text
#  name        :string           not null
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#
class StaffMember::Role < ApplicationRecord
  self.table_name = "staff_member_roles"

  has_many :permissions, class_name: "StaffMember::Role::Permission", dependent: :destroy
  has_many :staff_members, class_name: "StaffMember", dependent: :restrict_with_error

  validates :name, presence: true

  sig { params(_auth_object: T.untyped).returns(T::Array[String]) }
  def self.ransackable_attributes(_auth_object = nil)
    %w[id name created_at]
  end
end
