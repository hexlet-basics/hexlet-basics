# frozen_string_literal: true

class Web::Admin::ApplicationController < Web::ApplicationController
  before_action :authenticate_admin!

  before_action do
    seo_tags = {
      title: t(".title")
    }
    set_meta_tags seo_tags
  end
end
