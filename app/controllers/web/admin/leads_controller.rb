class Web::Admin::LeadsController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "id", "so" => "desc")
    search = Lead.includes([ :user ]).ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      leads: LeadResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end
end
