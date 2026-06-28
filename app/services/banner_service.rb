# typed: strict
# frozen_string_literal: true

class BannerService < ApplicationService
  class << self
    extend T::Sig

    sig { params(struct: BannerStruct).returns(Typed::Result[Banner, Banner]) }
    def create(struct)
      banner = Banner.new(struct.attributes)
      return fail_with(banner) unless banner.save

      success_with(banner)
    end

    sig { params(id: T.untyped, struct: BannerStruct).returns(Typed::Result[Banner, Banner]) }
    def update(id, struct)
      banner = Banner.find(id)
      return fail_with(banner) unless banner.update(struct.attributes)

      success_with(banner)
    end
  end
end
