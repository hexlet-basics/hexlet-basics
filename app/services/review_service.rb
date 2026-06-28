# typed: strict
# frozen_string_literal: true

class ReviewService < ApplicationService
  class << self
    extend T::Sig

    sig { params(struct: ReviewStruct, locale: String).returns(Typed::Result[Review, Review]) }
    def create(struct, locale:)
      review = Review.new(struct.attributes.merge(locale:))
      return fail_with(review) unless review.save

      success_with(review)
    end

    sig { params(id: T.untyped, struct: ReviewStruct, locale: String).returns(Typed::Result[Review, Review]) }
    def update(id, struct, locale:)
      review = Review.find(id)
      review.assign_attributes(struct.attributes.merge(locale:))
      return fail_with(review) unless review.save

      success_with(review)
    end
  end
end
