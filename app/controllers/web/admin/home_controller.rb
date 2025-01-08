# frozen_string_literal: true

class Web::Admin::HomeController < Web::Admin::ApplicationController
  def index
    q = params.fetch(:q, {}).with_defaults("s" => "created_at desc")
    search = User.admin.ransack(q)
    users = search.result

    render inertia: true, props: {
      users:
    }
  end
end
