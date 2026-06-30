# typed: strict
# frozen_string_literal: true

class Web::Admin::Management::RolesController < Web::Admin::Management::ApplicationController
  sig { void }
  def index
    roles = StaffMember::Role.includes(:permissions).order(:name)

    render inertia: true, props: {
      roles: StaffRoleResource.new(roles),
      totalPermissionsCount: StaffMember::Role::Permission::Resource.values.count
    }
  end

  sig { void }
  def new
    role = StaffMember::Role.new

    render inertia: true, props: {
      roleCrud: StaffRoleCrudResource.new(role)
    }
  end

  sig { void }
  def edit
    role = StaffMember::Role.find(params[:id])

    render inertia: true, props: {
      roleCrud: StaffRoleCrudResource.new(role)
    }
  end

  sig { void }
  def create
    struct = ApplicationParamsStruct.from_params!(StaffRoleCreateStruct, params.require(:data))
    result = StaffRoleService.create(struct)

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_admin_management_role_url(result.payload.role)
    when Typed::Failure
      f(:error)
      redirect_to new_admin_management_role_url, inertia: { errors: result.error.role.errors }
    end
  end

  sig { void }
  def update
    struct = ApplicationParamsStruct.from_params!(StaffRoleUpdateStruct, params.require(:data))
    result = StaffRoleService.update(params[:id], struct)

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_admin_management_role_url(result.payload.role)
    when Typed::Failure
      f(:error)
      redirect_to edit_admin_management_role_url(result.error.role), inertia: { errors: result.error.role.errors }
    end
  end

  sig { void }
  def destroy
    role = StaffMember::Role.find(params[:id])
    role.destroy ? f(:success) : f(:error)
    redirect_to admin_management_roles_url
  end
end
