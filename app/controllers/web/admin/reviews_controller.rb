# frozen_string_literal: true

class Web::Admin::ReviewsController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "created_at", "so" => "0")
    search = Review.includes([ :user ]).ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      reviews: ReviewResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  def new
    review = Admin::ReviewForm.new
    languages = Language.all

    render inertia: true, props: {
      reviewDto: ReviewCrudResource.new(review),
      courses: LanguageResource.new(languages)
    }
  end

  def edit
    review = Admin::ReviewForm.find(params[:id])
    languages = Language.all

    render inertia: true, props: {
      reviewDto: ReviewCrudResource.new(review),
      courses: LanguageResource.new(languages)
    }
  end

  def create
    review = Admin::ReviewForm.new(params[:review])

    if review.save
      f(:success)
      redirect_to edit_admin_review_url(review)
    else
      f(:error)
      redirect_to_inertia new_admin_review_url, review
    end
  end

  def update
    review = Admin::ReviewForm.find(params[:id])

    if review.update(params[:review])
      f(:success)
    else
      f(:error)
    end

    redirect_to_inertia edit_admin_review_url(review), review
  end
end
