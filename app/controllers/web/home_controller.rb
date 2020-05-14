# frozen_string_literal: true

class Web::HomeController < Web::ApplicationController
  def index
    @languages = Language.all
  end
end
