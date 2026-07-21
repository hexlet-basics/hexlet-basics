# typed: strict
# frozen_string_literal: true

class ReviewService < ApplicationService
  class << self
    extend T::Sig

    sig { params(struct: ReviewStruct, locale: String).returns(Typed::Result[Review, Review]) }
    def create(struct, locale:)
      attributes = struct.attributes.merge(locale:)
      review = Review.new(attributes)
      return fail_with(review) unless review.save

      success_with(review)
    end

    sig { params(id: String, struct: ReviewStruct, locale: String).returns(Typed::Result[Review, Review]) }
    def update(id, struct, locale:)
      review = Review.find(id)
      attributes = struct.attributes.merge(locale:)
      return fail_with(review) unless review.update(attributes)

      success_with(review)
    end
  end
end
