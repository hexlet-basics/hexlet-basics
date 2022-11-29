# frozen_string_literal: true

class Web::BlogPostsController < Web::ApplicationController
  def index
    @blog_posts = BlogPost.published.page(params[:page])
  end

  def show
    @blog_post = BlogPost.published.find_by! slug: (params[:id])
  end
end
