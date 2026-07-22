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

    sig { params(id: String, struct: LanguageStruct, cover: T.untyped).returns(Typed::Result[Language, Language]) }
    def update(id, struct, cover: nil)
      language = Language.find(id)
      attributes = attrs(struct, cover)
      return fail_with(language) unless language.update(attributes)

      success_with(language)
    end

    private

    # cover lives outside the struct; set it only for a real upload — forms may
    # echo back a serialized attachment hash, which is not attachable
    sig { params(struct: LanguageStruct, cover: T.untyped).returns(T::Hash[Symbol, T.untyped]) }
    def attrs(struct, cover)
      result = struct.attributes
      result[:cover] = cover if cover.is_a?(ActionDispatch::Http::UploadedFile)
      result
    end
  end
end
