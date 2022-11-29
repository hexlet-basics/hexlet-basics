# frozen_string_literal: true

class Web::Admin::ReviewsController < Web::Admin::ApplicationController
  def index
    q = params.fetch(:q, {}).with_defaults('s' => 'created_at desc')
    @search = Review.ransack(q)
    @reviews = @search.result
  end

  def new
    @review = Review.new
  end

  def edit
    @review = Review.find(params[:id])
  end

  def create
    @review = Review.new(review_params)

    if @review.save
      f(:success)
      redirect_to admin_reviews_path
    else
      f(:error)
      render :new
    end
  end

  def update
    @review = Review.find(params[:id])

    if @review.update(review_params)
      f(:success)
      redirect_to admin_reviews_path
    else
      f(:error)
      render :edit
    end
  end

  private

  def review_params
    params.require(:review).permit(:user_id, :state_event, :language_id, :locale, :body)
  end
end
