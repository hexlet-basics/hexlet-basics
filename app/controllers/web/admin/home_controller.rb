# typed: strict
# frozen_string_literal: true

class Web::Admin::HomeController < Web::Admin::ApplicationController
  # Дашборд доступен любому staff; список админов отдаём только суперпользователю.
  STAFF_RESOURCE = nil

  sig { returns(T.untyped) }
  def index
    q = ransack_params("sf" => "id", "so" => "desc")
    search = User.admin.ransack(q)
    admins = current_user&.admin? ? search.result : User.none

    render inertia: true, props: {
      admins: UserResource.new(admins),
      grid: GridResource.new(grid_params)
    }
  end
end
