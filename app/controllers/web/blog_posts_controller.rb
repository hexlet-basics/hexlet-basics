class Web::BlogPostsController < Web::ApplicationController
  def index
    scope = BlogPost.published.with_locale.includes([ :cover_attachment ])
    pagy, records = pagy(scope)

    seo_tags = {
      title: t(".header"),
      description: t(".meta.description"),
      canonical: blog_posts_url
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      blogPosts: BlogPostResource.new(records),
      pagy:
    }
  end

  def show
    blog_post = BlogPost.published.find_by!(slug: params[:id])

    # category = blog_post.category

    # blog_posts = []
    # languages = []

    blog_posts = BlogPost.published.with_locale.except(blog_post).includes([ :cover_attachment ]).limit(2)

    # if category
    #   blog_posts = category.blog_posts.except(blog_post).limit(3)
    #   languages = category.languages.limit(3)
    # end

    seo_tags = {
      title: blog_post.name,
      canonical: blog_post_url(blog_post.slug)
    }
    set_meta_tags seo_tags

    # TODO: add https://developers.google.com/search/docs/appearance/structured-data/article
    render inertia: true, props: {
      blogPost: BlogPostResource.new(blog_post),
      recommendedBlogPosts: BlogPostResource.new(blog_posts)
    }
  end
end
