# frozen_string_literal: true

class Web::ReviewsController < Web::ApplicationController
  def index
    @reviews = Review.published.with_locale.order(id: :desc).page(params[:page])
  end
end
