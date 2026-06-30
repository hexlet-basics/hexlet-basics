# typed: strict

class Web::Admin::ApplicationController < Web::ApplicationController
  before_action :authenticate_staff!
  before_action :authorize_staff_resource!

  before_action do
    T.bind(self, Web::Admin::ApplicationController)
    seo_tags = {
      title: t(".title")
    }
    set_meta_tags seo_tags
  end

  inertia_share do
    T.bind(self, Web::Admin::ApplicationController)
    staff = current_user&.staff_member
    {
      isAdmin: current_user&.admin? || false,
      staffPermissions: staff ? staff.role.permissions.index_by(&:resource).transform_values { |p|
        p.slice(:can_index, :can_create, :can_update, :can_destroy)
      } : {}
    }
  end

  private

  sig { void }
  def authenticate_staff!
    require_authentication
    redirect_to root_path unless current_user&.staff?
  end

  # STAFF_RESOURCE семантика:
  #   константа = Resource::X → проверяем право по action
  #   = nil                   → доступ любому staff (напр. дашборд)
  #   не объявлена            → только admin
  sig { void }
  def authorize_staff_resource!
    return if current_user&.admin?

    unless self.class.const_defined?(:STAFF_RESOURCE, false)
      redirect_to admin_root_path
      return
    end

    resource = T.cast(
      self.class.const_get(:STAFF_RESOURCE, false),
      T.nilable(StaffMember::Role::Permission::Resource)
    )
    return if resource.nil?

    staff = current_user&.staff_member
    serialized = resource.serialize
    allowed = case action_name
    when "index", "show" then staff&.can_index_resource?(serialized)
    when "new", "create" then staff&.can_create_resource?(serialized)
    when "edit", "update", "review", "related_courses" then staff&.can_update_resource?(serialized)
    when "destroy" then staff&.can_destroy_resource?(serialized)
    else false
    end

    redirect_to admin_root_path unless allowed
  end
end
