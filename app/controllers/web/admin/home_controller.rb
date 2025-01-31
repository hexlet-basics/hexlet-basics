# frozen_string_literal: true

class Web::Admin::HomeController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "created_at", "so" => "0")
    search = User.admin.ransack(q)
    admins = search.result

    render inertia: true, props: {
      admins: UserResource.new(admins),
      grid: GridResource.new(grid_params)
    }
  end
end
