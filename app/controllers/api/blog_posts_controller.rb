# typed: strict

class Api::BlogPostsController < Api::ApplicationController
  allow_unauthenticated_access

  sig { returns(T.untyped) }
  def show
    blog_post = BlogPost.find(params[:id])
    respond_with BlogPostResource.new(blog_post)
  end

  sig { returns(T.untyped) }
  def next
    next_blog_post = BlogPost.with_locale.published_state.where("id < ?", params[:id]).order(id: :desc).first!
    respond_with BlogPostResource.new(next_blog_post)
  end
end
