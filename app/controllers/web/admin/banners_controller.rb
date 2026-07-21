# typed: strict
# frozen_string_literal: true

class Web::Admin::BannersController < Web::Admin::ApplicationController
  STAFF_RESOURCE = StaffMember::Role::Permission::Resource::Banners

  sig { void }
  def index
    q = ransack_params("sf" => "id", "so" => "desc")
    search = Banner.ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      banners: BannerResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  sig { void }
  def new
    banner = Banner.new

    render inertia: true, props: {
      bannerDto: BannerCreateResource.new(banner)
    }
  end

  sig { void }
  def edit
    banner = Banner.find(params[:id])

    render inertia: true, props: {
      bannerDto: BannerUpdateResource.new(banner)
    }
  end

  sig { void }
  def create
    struct = ApplicationParamsStruct.from_params(BannerStruct, params.require(:data))
    result = BannerService.create(struct)

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_admin_banner_url(result.payload)
    else
      f(:error)
      redirect_to new_admin_banner_url, inertia: { errors: result.error.errors }
    end
  end

  sig { void }
  def update
    struct = ApplicationParamsStruct.from_params(BannerStruct, params.require(:data))
    result = BannerService.update(params[:id], struct)

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_admin_banner_url(result.payload)
    when Typed::Failure
      f(:error)
      redirect_to edit_admin_banner_url(result.error), inertia: { errors: result.error.errors }
    end
  end

  sig { void }
  def destroy
    banner = Banner.find(params[:id])
    banner.destroy ? f(:success) : f(:error)

    redirect_to admin_banners_url
  end
end
