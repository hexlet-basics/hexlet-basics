# frozen_string_literal: true

class Web::ReviewsController < Web::ApplicationController
  def index
    scope = Review.published.with_locale.includes([ :user ]).order(id: :desc)
    pagy, records = pagy(scope)

    seo_tags = {
      title: t(".title"),
      canonical: reviews_url
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      reviews: ReviewResource.new(records),
      pagy:
    }
  end
end
