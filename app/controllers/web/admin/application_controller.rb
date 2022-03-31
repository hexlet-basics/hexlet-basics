# frozen_string_literal: true

class Web::Admin::ApplicationController < Web::ApplicationController
  before_action :authenticate_admin!

  before_action do
    set_meta_tags title: t('title', scope: 'web.admin')
  end
end
