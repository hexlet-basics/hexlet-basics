# frozen_string_literal: true

class Web::Admin::BlogPostsController < Web::Admin::ApplicationController
  def index
    q = params.fetch(:q, {}).with_defaults('s' => 'created_at desc')
    @search = BlogPost.ransack(q)
    @blog_posts = @search.result
  end

  def new
    @blog_post = Admin::BlogPostForm.new
    @blog_post.creator = current_user
  end

  def edit
    @blog_post = Admin::BlogPostForm.find(params[:id])
  end

  def create
    @blog_post = Admin::BlogPostForm.new(params[:admin_blog_post_form])
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
    @blog_post = Admin::BlogPostForm.find(params[:id])

    if @blog_post.update(params[:admin_blog_post_form])
      f(:success)
      redirect_to admin_blog_posts_path
    else
      f(:error)
      render :edit
    end
  end
end
