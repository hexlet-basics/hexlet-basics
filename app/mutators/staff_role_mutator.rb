# typed: strict
# frozen_string_literal: true

class StaffRoleMutator
  class << self
    extend T::Sig

    sig { params(struct: StaffRoleCreateStruct).returns(StaffMember::Role) }
    def create(struct)
      role = StaffMember::Role.new(**struct.attributes)
      role.save
      role
    end

    sig { params(role: StaffMember::Role, struct: StaffRoleUpdateStruct).returns(T::Boolean) }
    def update(role, struct)
      role.update(struct.attributes)
    end

    sig { params(role: StaffMember::Role, permissions_params: T::Array[ActionController::Parameters]).void }
    def sync_permissions(role, permissions_params)
      permissions_params.each do |perm|
        resource = perm[:resource].to_s
        next if resource.blank?

        role.permissions.find_or_initialize_by(resource:).update!(
          can_index: perm[:can_index] == true || perm[:can_index] == "true",
          can_create: perm[:can_create] == true || perm[:can_create] == "true",
          can_update: perm[:can_update] == true || perm[:can_update] == "true",
          can_destroy: perm[:can_destroy] == true || perm[:can_destroy] == "true"
        )
      end
    end
  end
end
