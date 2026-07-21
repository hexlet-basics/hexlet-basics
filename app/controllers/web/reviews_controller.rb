# typed: strict
# frozen_string_literal: true

class Web::ReviewsController < Web::ApplicationController
  allow_unauthenticated_access

  sig { returns(T.untyped) }
  def index
    scope = Review.published_state.with_locale
      .includes([ :user, :language ])
      .order(id: :desc)
    pagy, records = T.unsafe(self).pagy(scope)

    description = t(".meta.description").truncate(160)
    seo_tags = {
      title: t(".header"),
      description: description,
      canonical: reviews_url,
      og: {
        title: t(".title"),
        description:  description
      },
      twitter: {
        card: "summary",
        site: t("links.hexlet_twitter_handle")
      }
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      reviews: ReviewResource.new(records),
      pagy: PagyResource.new(pagy)
    }
  end
end
