# frozen_string_literal: true

class Web::ReviewsController < Web::ApplicationController
  def index
    @reviews = Review.published.page(params[:page])
  end

  def show
    @review = Review.published.find params[:id]
  end
end
