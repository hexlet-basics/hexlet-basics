# frozen_string_literal: true

class Web::HomeController < Web::ApplicationController
  def index
    @languages = Language.joins(:current_version).all
  end
end
