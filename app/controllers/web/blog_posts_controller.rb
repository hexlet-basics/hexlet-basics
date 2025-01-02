# frozen_string_literal: true

class Web::BlogPostsController < Web::ApplicationController
  def index
    scope = BlogPost.published.with_locale
    pagy, records = pagy(scope)

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

    blog_posts = BlogPost.with_locale.except(blog_post).limit(2)

    # if category
    #   blog_posts = category.blog_posts.except(blog_post).limit(3)
    #   languages = category.languages.limit(3)
    # end

    # TODO: add https://developers.google.com/search/docs/appearance/structured-data/article
    render inertia: true, props: {
      blogPost: BlogPostResource.new(blog_post),
      recommendedBlogPosts: BlogPostResource.new(blog_posts)
    }
  end
end
