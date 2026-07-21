# typed: strict

class Web::Admin::BlogPostsController < Web::Admin::ApplicationController
  STAFF_RESOURCE = StaffMember::Role::Permission::Resource::BlogPosts

  sig { void }
  def index
    default_params = { "sf" => "id", "so" => "desc" }
    q = ransack_params(default_params)
    search = BlogPost.with_locale.includes([ :likes, :cover_attachment ]).ransack(q)
    pagy, records = pagy(search.result)
    # raise q.inspect

    grid = GridResource.new(grid_params(pagy, default_params))
    blog_posts = BlogPostResource.new(records)

    render inertia: true, props: {
      blogPosts: blog_posts,
      grid:
    }
  end

  sig { void }
  def new
    blog_post = BlogPost.new
    blog_post.creator = T.must(current_user)

    render inertia: true, props: {
      blogPostDto: BlogPostCreateResource.new(blog_post)
    }
  end

  sig { void }
  def edit
    blog_post = BlogPost.find(params[:id])

    render inertia: true, props: {
      blogPostDto: BlogPostUpdateResource.new(blog_post),
      relatedCourses: LanguageResource.new(blog_post.related_languages)
    }
  end

  sig { void }
  def related_courses
    FindRelatedCoursesForBlogPostJob.perform_later(params[:id])

    f(:success)
    redirect_back_or_to admin_blog_posts_path
  end

  sig { void }
  def create
    struct = ApplicationParamsStruct.from_params(BlogPostStruct, params.require(:data))
    result = BlogPostService.create(struct, creator: T.must(current_user), locale: I18n.locale.to_s, cover: params.dig(:data, :cover))

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_admin_blog_post_path(result.payload)
    when Typed::Failure
      f(:error)
      redirect_to new_admin_blog_post_url, inertia: { errors: result.error.errors }
    end
  end

  sig { void }
  def update
    struct = ApplicationParamsStruct.from_params(BlogPostStruct, params.require(:data))
    result = BlogPostService.update(params[:id], struct, locale: I18n.locale.to_s, cover: params.dig(:data, :cover))

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_admin_blog_post_path(result.payload)
    when Typed::Failure
      f(:error)
      redirect_to edit_admin_blog_post_path(result.error), inertia: { errors: result.error.errors }
    end
  end
end
