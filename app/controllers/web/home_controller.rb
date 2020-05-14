# frozen_string_literal: true

class Web::HomeController < ApplicationController
  def index
    @languages = Language.all
  end
end
