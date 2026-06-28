# typed: strict
# frozen_string_literal: true

class BlogPostService < ApplicationService
  class << self
    extend T::Sig

    sig { params(struct: BlogPostStruct, creator: User, locale: String, cover: T.untyped).returns(Typed::Result[BlogPost, BlogPost]) }
    def create(struct, creator:, locale:, cover: nil)
      blog_post = BlogPost.new(attrs(struct, cover).merge(creator:, locale:))
      return fail_with(blog_post) unless blog_post.save

      success_with(blog_post)
    end

    sig { params(id: T.untyped, struct: BlogPostStruct, locale: String, cover: T.untyped).returns(Typed::Result[BlogPost, BlogPost]) }
    def update(id, struct, locale:, cover: nil)
      blog_post = BlogPost.find(id)
      blog_post.assign_attributes(attrs(struct, cover).merge(locale:))
      return fail_with(blog_post) unless blog_post.save

      success_with(blog_post)
    end

    private

    # the cover file isn't part of the typed struct; only set it when present so
    # an empty file input doesn't clear the existing attachment
    sig { params(struct: BlogPostStruct, cover: T.untyped).returns(T::Hash[Symbol, T.untyped]) }
    def attrs(struct, cover)
      result = struct.attributes
      result[:cover] = cover if cover.present?
      result
    end
  end
end
