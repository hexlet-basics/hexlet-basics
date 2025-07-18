class Api::BlogPostsController < Api::ApplicationController
  def show
    blog_post = BlogPost.find(params[:id])
    render json: BlogPostResource.new(blog_post)
  end

  def next
    next_blog_post = BlogPost.with_locale.published_state.where("id < ?", params[:id]).order(id: :desc).first!
    render json: BlogPostResource.new(next_blog_post)
  end
end
