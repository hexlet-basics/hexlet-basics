# typed: strict
# frozen_string_literal: true

class Web::Admin::Management::RolePermissionsController < Web::Admin::Management::ApplicationController
  sig { void }
  def show
    role = StaffMember::Role.includes(:permissions).find(params.expect(:id))

    render inertia: true, props: {
      role: StaffRoleCrudResource.new(role),
      resources: StaffMember::Role::Permission::Resource.values.map(&:serialize)
    }
  end

  sig { void }
  def update
    role = StaffMember::Role.find(params.expect(:id))
    permissions_params = params.require(:data).require(:permissions)
    StaffRoleMutator.sync_permissions(role, permissions_params)

    f(:success)
    redirect_to admin_management_role_permission_url(role)
  end
end
