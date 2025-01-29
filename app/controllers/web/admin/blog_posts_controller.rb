# frozen_string_literal: true

class Web::Admin::BlogPostsController < Web::Admin::ApplicationController
  def index
    q = ransack_params("s" => "created_at desc")
    search = BlogPost.includes([ :cover_attachment ]).ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      blogPosts: BlogPostResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  def new
    blog_post = Admin::BlogPostForm.new
    blog_post.creator = current_user

    render inertia: true, props: {
      blog_post: BlogPostResource.new(blog_post)
    }
  end

  def edit
    blog_post = Admin::BlogPostForm.find(params[:id])

    render inertia: true, props: {
      blog_post: BlogPostResource.new(blog_post)
    }
  end

  def create
    blog_post = Admin::BlogPostForm.new(params[:blog_post])
    blog_post.creator = current_user

    if blog_post.save
      f(:success)
      redirect_to_inertia edit_admin_blog_post_path(blog_post), blog_post
    else
      f(:error)
      redirect_to_inertia new_admin_blog_post_url, blog_post
    end
  end

  def update
    blog_post = Admin::BlogPostForm.find(params[:id])
    # raise params[:cover].inspect

    if blog_post.update(params[:blog_post])
      f(:success)
    else
      f(:error)
    end

      redirect_to_inertia edit_admin_blog_post_path(blog_post), blog_post
  end
end
