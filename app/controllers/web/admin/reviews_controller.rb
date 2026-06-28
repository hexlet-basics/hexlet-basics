# typed: true
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
    review = Review.new
    languages = Language.all

    render inertia: true, props: {
      reviewDto: ReviewCreateResource.new(review),
      courses: LanguageResource.new(languages)
    }
  end

  def edit
    review = Review.find(params[:id])
    languages = Language.all

    render inertia: true, props: {
      reviewDto: ReviewUpdateResource.new(review),
      courses: LanguageResource.new(languages)
    }
  end

  def create
    struct = ApplicationParamsStruct.from_params(ReviewStruct, params.require(:data))
    result = ReviewService.create(struct, locale: I18n.locale.to_s)

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_admin_review_url(result.payload)
    else
      f(:error)
      redirect_to new_admin_review_url, inertia: { errors: result.error.errors }
    end
  end

  def update
    struct = ApplicationParamsStruct.from_params(ReviewStruct, params.require(:data))
    result = ReviewService.update(params[:id], struct, locale: I18n.locale.to_s)

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_admin_review_url(result.payload)
    when Typed::Failure
      f(:error)
      redirect_to edit_admin_review_url(result.error), inertia: { errors: result.error.errors }
    end
  end
end
