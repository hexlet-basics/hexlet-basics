# frozen_string_literal: true

class Web::Admin::ReviewsController < Web::Admin::ApplicationController
  def index
    @reviews = Language.page(params[:page])
  end

  def new; end

  def edit; end
end
