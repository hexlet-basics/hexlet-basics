# frozen_string_literal: true

class Web::Admin::BlogPostsController < Web::Admin::ApplicationController
  def index
    q = params.fetch(:q, {}).with_defaults('s' => 'created_at desc')
    @search = BlogPost.ransack(q)
    @blog_posts = @search.result
  end

  def new
    @blog_post = BlogPost.new
  end

  def edit
    @blog_post = BlogPost.find(params[:id])
  end

  def create
    @blog_post = BlogPost.new(blog_post_params)
    @blog_post.creator = current_user

    if @blog_post.save
      f(:success)
      redirect_to admin_blog_posts_path
    else
      f(:error)
      render :new
    end
  end

  def update
    @blog_post = BlogPost.find(params[:id])

    if @blog_post.update(blog_post_params)
      f(:success)
      redirect_to admin_blog_posts_path
    else
      f(:error)
      render :edit
    end
  end

  private

  def blog_post_params
    params.require(:blog_post).permit(:creator_id, :language_id, :body, :slug, :locale, :name)
  end
end
