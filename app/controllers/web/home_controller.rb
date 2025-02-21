class Web::HomeController < Web::ApplicationController
  def index

    language_member_resources = current_user.language_members.includes([ language: :current_version ])
      .map { |m| Language::MemberResource.new(m) }
    language_member_resources_by_language = language_member_resources.index_by { |r| r.object.language_id }

    blog_posts = BlogPost.published.with_locale.includes(:cover_attachment).last(3)
    # reviews = Review.random
    #
    # # TODO: need refactor it
    # scope = Language.includes(:current_version)
    # @languages_links_by_slug = scope.each_with_object({}) do |item, acc|
    #   acc[item.slug.to_sym] = view_context.link_to(item, language_path(item.slug))
    # end
    #

    # categories = Language::Category.all
    #
    # completed_languages = Language.with_progress(:completed).with_locale.ordered
    # infos = Language::Version::Info.with_locale.where(language: completed_languages).includes(:language)
    # infos_by_language = infos.index_by { |item| item.language.id }
    # item_builders = completed_languages.map { |l| CourseSchema.to_builder(l, infos_by_language.fetch(l.id)) }
    #
    # @builder = ItemListSchema.to_builder(item_builders)

    # gon.languages_for_widget = helpers.completed_languages.map(&:name)

    seo_tags = {
      title: t(".title"),
      description: t(".meta.description"),
      canonical: root_url,
      image_src: view_context.vite_asset_path("images/logo.png"),
      alternate: {
        ru: view_context.root_url(suffix: :ru),
        en: view_context.root_url(suffix: nil)
      },
      twitter: {
        card: "summary",
        site: "@hexlet_io"
      },
      og: {
        title: t(".title"),
        type: "website",
        url: root_url,
        image: view_context.vite_asset_path("images/logo.png")
      }
    }
    set_meta_tags seo_tags

    user = User::SignUpForm.new

    render inertia: true, props: {
      courseMembersByCourseId: language_member_resources_by_language,
      blogPosts: BlogPostResource.new(blog_posts),
      newUser: UserSignUpFormResource.new(user),
      courseCategories: Language::CategoryResource.new(Language::Category.all)
    }
  end

  def robots
    respond_to :text
  end
end
