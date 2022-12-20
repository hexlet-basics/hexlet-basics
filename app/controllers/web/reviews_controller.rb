# frozen_string_literal: true

class Web::ReviewsController < Web::ApplicationController
  def index
    @reviews = Review.published.with_locale.page(params[:page])
  end

  def show
    @review = Review.published.find params[:id]
    # @builder = ReviewSchema.to_builder(@review)
  end
end
