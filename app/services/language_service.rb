# typed: strict
# frozen_string_literal: true

class LanguageService < ApplicationService
  class << self
    extend T::Sig

    sig { params(struct: LanguageStruct, cover: T.untyped).returns(Typed::Result[Language, Language]) }
    def create(struct, cover: nil)
      attributes = attrs(struct, cover)
      language = Language.new(attributes)
      return fail_with(language) unless language.save

      success_with(language)
    end

    sig { params(id: T.untyped, struct: LanguageStruct, cover: T.untyped).returns(Typed::Result[Language, Language]) }
    def update(id, struct, cover: nil)
      language = Language.find(id)
      attributes = attrs(struct, cover)
      return fail_with(language) unless language.update(attributes)

      success_with(language)
    end

    private

    # cover lives outside the struct; set it only when present (skip_if_empty)
    sig { params(struct: LanguageStruct, cover: T.untyped).returns(T::Hash[Symbol, T.untyped]) }
    def attrs(struct, cover)
      result = struct.attributes
      result[:cover] = cover if cover.present?
      result
    end
  end
end
