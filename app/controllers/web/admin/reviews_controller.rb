# frozen_string_literal: true

class Web::Admin::ReviewsController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "id", "so" => "desc")
    search = Review.with_locale.includes([ :user ]).ransack(q)
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
      reviewDto: ReviewCreateResource.new(review),
      courses: LanguageResource.new(languages)
    }
  end

  def edit
    review = Admin::ReviewForm.find(params[:id])
    languages = Language.all

    render inertia: true, props: {
      reviewDto: ReviewUpdateResource.new(review),
      courses: LanguageResource.new(languages)
    }
  end

  def create
    review = Admin::ReviewForm.new(params[:data])
    review.locale = I18n.locale

    if review.save
      f(:success)
      redirect_to edit_admin_review_url(review)
    else
      f(:error)
      redirect_to new_admin_review_url, inertia: { errors: review.errors }
    end
  end

  def update
    review = Admin::ReviewForm.find(params[:id])
    review.locale = I18n.locale

    if review.update(params[:data])
      f(:success)
    else
      f(:error)
    end

    redirect_to edit_admin_review_url(review), inertia: { errors: review.errors }
  end
end
