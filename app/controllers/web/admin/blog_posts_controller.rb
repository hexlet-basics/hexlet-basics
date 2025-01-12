# frozen_string_literal: true

class Web::Admin::BlogPostsController < Web::Admin::ApplicationController
  def index
    q = ransack_params("s" => "created_at desc")
    search = BlogPost.ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      blogPosts: BlogPostResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  def new
    blog_post = Admin::BlogPostForm.new
    blog_post.creator = current_user

    render inertia: true, props: {}
  end

  def edit
    blog_post = Admin::BlogPostForm.find(params[:id])

    render inertia: true, props: {
      blog_post: BlogPostResource.new(blog_post)
    }
  end

  def create
    blog_post = Admin::BlogPostForm.new(params[:admin_blog_post_form])
    blog_post.creator = current_user

    if blog_post.save
      f(:success)
      redirect_to admin_blog_posts_path
    else
      f(:error)
      redirect_to_inertia edit_admin_blog_post_path(user), user
    end
  end

  def update
    blog_post = Admin::BlogPostForm.find(params[:id])

    if blog_post.update(params[:admin_blog_post_form])
      f(:success)
      redirect_to admin_blog_posts_path
    else
      f(:error)
      redirect_to_inertia edit_admin_blog_post_path(user), user
    end
  end
end
