# frozen_string_literal: true

class Web::BlogPostsController < Web::ApplicationController
  def index
    @blog_posts = BlogPost.published.with_locale.page(params[:page])
  end

  def show
    @blog_post = BlogPost.published.find_by!(slug: params[:id])

    @category = @blog_post.category
    @blog_posts = []
    @languages = []
    if @category
      @blog_posts = @category.blog_posts.except(@blog_post).limit(3)
      @languages = @category.languages.limit(3)
    end

    # TODO: add https://developers.google.com/search/docs/appearance/structured-data/article
  end
end
