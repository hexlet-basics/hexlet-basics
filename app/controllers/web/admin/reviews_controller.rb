# frozen_string_literal: true

class Web::Admin::ReviewsController < Web::Admin::ApplicationController
  def index
    q = params.fetch(:q, {}).with_defaults('s' => 'created_at desc')
    @search = Review.ransack(q)
    @reviews = @search.result
  end

  def new
    @review = Admin::ReviewForm.new
  end

  def edit
    @review = Admin::ReviewForm.find(params[:id])
  end

  def create
    @review = Admin::ReviewForm.new(params[:admin_review_form])

    if @review.save
      f(:success)
      redirect_to admin_reviews_path
    else
      f(:error)
      render :new
    end
  end

  def update
    @review = Admin::ReviewForm.find(params[:id])

    if @review.update(params[:admin_review_form])
      f(:success)
      redirect_to admin_reviews_path
    else
      f(:error)
      render :edit
    end
  end
end
