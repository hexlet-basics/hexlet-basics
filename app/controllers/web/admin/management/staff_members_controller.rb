# typed: strict
# frozen_string_literal: true

class Web::Admin::Management::StaffMembersController < Web::Admin::Management::ApplicationController
  sig { void }
  def index
    q = ransack_params("sf" => "id", "so" => "desc")
    search = StaffMember.includes(:user, :role).ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      staffMembers: StaffMemberResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  sig { void }
  def new
    staff_member = StaffMember.new

    render inertia: true, props: {
      staffMemberCrud: StaffMemberCrudResource.new(staff_member),
      roles: StaffRoleResource.new(StaffMember::Role.order(:name)),
      locales: I18n.available_locales.map(&:to_s)
    }
  end

  sig { void }
  def edit
    staff_member = StaffMember.includes(:role, :user).find(params[:id])

    render inertia: true, props: {
      staffMemberCrud: StaffMemberCrudResource.new(staff_member),
      roles: StaffRoleResource.new(StaffMember::Role.order(:name)),
      locales: I18n.available_locales.map(&:to_s)
    }
  end

  sig { void }
  def create
    struct = ApplicationParamsStruct.from_params!(StaffMemberCreateStruct, params.require(:data))
    result = StaffMemberService.create(struct)

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_admin_management_staff_member_url(result.payload.staff_member)
    when Typed::Failure
      f(:error)
      redirect_to new_admin_management_staff_member_url, inertia: { errors: result.error.staff_member.errors }
    end
  end

  sig { void }
  def update
    struct = ApplicationParamsStruct.from_params!(StaffMemberUpdateStruct, params.require(:data))
    result = StaffMemberService.update(params[:id], struct)

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_admin_management_staff_member_url(result.payload.staff_member)
    when Typed::Failure
      f(:error)
      redirect_to edit_admin_management_staff_member_url(result.error.staff_member), inertia: { errors: result.error.staff_member.errors }
    end
  end

  sig { void }
  def destroy
    staff_member = StaffMember.find(params[:id])
    staff_member.destroy ? f(:success) : f(:error)
    redirect_to admin_management_staff_members_url
  end
end
