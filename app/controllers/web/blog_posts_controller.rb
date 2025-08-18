class Web::BlogPostsController < Web::ApplicationController
  def index
    scope = BlogPost.published_state.with_locale
      .includes([ :creator, { cover_attachment: :blob } ])
      .order(id: :desc)
    pagy, records = pagy(scope)

    seo_tags = {
      title: t(".header"),
      description: t(".meta.description"),
      canonical: blog_posts_url,
      og: {
        title: t(".header"),
        description: t(".meta.description")
      },
      twitter: {
        card: "summary",
        site: "@hexlethq"
      }
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      blogPosts: BlogPostResource.new(records),
      pagy:
    }
  end

  def show
    blog_post = BlogPost.with_locale.published_state.find_by!(slug: params[:id])

    blog_posts = BlogPost.published_state.with_locale
      .includes([ :creator, { cover_attachment: :blob } ])
      .except(blog_post)
      .limit(2)

    related_landings = blog_post.related_main_language_landing_pages.merge(BlogPost::RelatedLanguageItem.order(order: :asc))

    image_url = blog_post.cover.attached? && view_context.rails_representation_url(blog_post.cover.variant(:main))
    seo_tags = {
      title: blog_post.name,
      description: blog_post.description,
      canonical: blog_post_url(blog_post.slug),
      twitter: {
        card: "summary",
        site: "@hexlethq"
      },
      og: {
        title: blog_post.name,
        description: blog_post.description,
        image: image_url
      }
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      blogPost: BlogPostResource.new(blog_post),
      recommendedBlogPosts: BlogPostResource.new(blog_posts),
      relatedLandingPages: Language::LandingPageForListsResource.new(related_landings)
    }
  end
end
