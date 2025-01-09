# frozen_string_literal: true

class Web::Admin::HomeController < Web::Admin::ApplicationController
  def index
    sp = search_params
    rp = sp[:ransack].with_defaults("s" => "created_at desc")
    search = User.admin.ransack(rp)
    admins = search.result

    render inertia: true, props: {
      admins: UserResource.new(admins),
      q: QResource.new(sp[:raw])
    }
  end
end
