class Web::Admin::BlogPostsController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "id", "so" => "desc")
    search = BlogPost.with_locale.includes([ :cover_attachment ]).ransack(q)
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
      blogPostDto: BlogPostCrudResource.new(blog_post)
    }
  end

  def edit
    blog_post = Admin::BlogPostForm.find(params[:id])

    render inertia: true, props: {
      blogPostDto: BlogPostCrudResource.new(blog_post)
    }
  end

  def related_courses
    FindRelatedCoursesForBlogPostJob.perform_later(params[:id])

    f(:success)
    redirect_back_or_to admin_blog_posts_path
  end

  def create
    blog_post = Admin::BlogPostForm.new(params[:blog_post])
    blog_post.locale = I18n.locale
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
    blog_post.locale = I18n.locale

    if blog_post.update(params[:blog_post])
      f(:success)
    else
      f(:error)
    end

      redirect_to_inertia edit_admin_blog_post_path(blog_post), blog_post
  end
end
