# typed: true

class Web::Admin::BlogPostsController < Web::Admin::ApplicationController
  def index
    default_params = { "sf" => "id", "so" => "desc" }
    q = ransack_params(default_params)
    search = BlogPost.with_locale.includes([ :cover_attachment ]).ransack(q)
    pagy, records = pagy(search.result)
    # raise q.inspect

    grid = GridResource.new(grid_params(pagy, default_params))
    blog_posts = BlogPostResource.new(records)

    render inertia: true, props: {
      blogPosts: blog_posts,
      grid:
    }
  end

  def new
    blog_post = Admin::BlogPostForm.new
    blog_post.creator = current_user

    render inertia: true, props: {
      blogPostDto: BlogPostCreateResource.new(blog_post)
    }
  end

  def edit
    blog_post = Admin::BlogPostForm.find(params[:id])

    render inertia: true, props: {
      blogPostDto: BlogPostUpdateResource.new(blog_post),
      relatedCourses: LanguageResource.new(blog_post.related_languages)
    }
  end

  def related_courses
    FindRelatedCoursesForBlogPostJob.perform_later(params[:id])

    f(:success)
    redirect_back_or_to admin_blog_posts_path
  end

  def create
    blog_post = Admin::BlogPostForm.new(params[:data])
    blog_post.locale = I18n.locale.to_s
    blog_post.creator = current_user

    if blog_post.save
      f(:success)
      redirect_to edit_admin_blog_post_path(blog_post), inertia: { errors: blog_post.errors }
    else
      f(:error)
      redirect_to new_admin_blog_post_url, inertia: { errors: blog_post.errors }
    end
  end

  def update
    blog_post = Admin::BlogPostForm.find(params[:id])
    blog_post.locale = I18n.locale.to_s

    if blog_post.update(params[:data])
      f(:success)
    else
      f(:error)
    end

      redirect_to edit_admin_blog_post_path(blog_post), inertia: { errors: blog_post.errors }
  end
end
