# typed: strict
# frozen_string_literal: true

class LandingPageService < ApplicationService
  LandingPage = Language::LandingPage

  class << self
    extend T::Sig

    sig { params(struct: LandingPageStruct, locale: String, outcomes_image: T.untyped).returns(Typed::Result[LandingPage, LandingPage]) }
    def create(struct, locale:, outcomes_image: nil)
      attributes = attrs(struct, outcomes_image).merge(locale:)
      landing_page = LandingPage.new(attributes)
      return fail_with(landing_page) unless landing_page.save

      success_with(landing_page)
    end

    sig { params(id: String, struct: LandingPageStruct, outcomes_image: T.untyped).returns(Typed::Result[LandingPage, LandingPage]) }
    def update(id, struct, outcomes_image: nil)
      landing_page = LandingPage.with_locale.find(id)
      attributes = attrs(struct, outcomes_image)
      return fail_with(landing_page) unless landing_page.update(attributes)

      success_with(landing_page)
    end

    private

    # outcomes_image lives outside the struct; set it only when present (skip_if_empty)
    sig { params(struct: LandingPageStruct, outcomes_image: T.untyped).returns(T::Hash[Symbol, T.untyped]) }
    def attrs(struct, outcomes_image)
      result = struct.attributes
      result[:outcomes_image] = outcomes_image if outcomes_image.present?
      result
    end
  end
end
