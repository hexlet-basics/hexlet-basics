# frozen_string_literal: true

class Web::Admin::ReviewsController < Web::Admin::ApplicationController
  def index
    q = ransack_params("s" => "created_at desc")
    search = Review.ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      reviews: ReviewResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  def new
    @review = Admin::ReviewForm.new

    render inertia: true, props: {
    }
  end

  def edit
    @review = Admin::ReviewForm.find(params[:id])

    render inertia: true, props: {
    }
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
