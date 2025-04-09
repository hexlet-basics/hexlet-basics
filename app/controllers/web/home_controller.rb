class Web::HomeController < Web::ApplicationController
  def index
    language_member_resources = current_user.language_members.includes([ language: :current_version ])
      .map { Language::MemberResource.new(it) }
    language_member_resources_by_language = language_member_resources.index_by { it.object.language_id }

    blog_posts = BlogPost.published
      .includes([ :creator, { cover_attachment: :blob } ])
      .with_locale
      .includes(:cover_attachment)
      .last(3)
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
      newUser: UserSignUpFormResource.new(user)
    }
  end

  def robots
    respond_to :text
  end

  def sitemap
    # NOTE: в запросах используется точечная выборка с помощью select из-за большого количества данных для уменьшения нагрузки
    # NOTE: данные выбираются специально по всем локалям, так как для карты сайта выводим ссылки для всех локалей.

    # NOTE: исключаем es, так как на сайте пока не работают ссылки связанный с es локалью
    locales = I18n.available_locales - [ I18n.locale, :es ]
    ordered_locales = [ I18n.locale, *locales ]

    language_landing_pages = Language::LandingPage
      .published.where(listed: true)
      .with_locale(ordered_locales)
      .includes(:language)
      .select(:language_id, :locale, :header, :slug, :id)

    language_landing_page_resources_by_locale = language_landing_pages
      .in_order_of(:locale, ordered_locales)
      .order(id: :asc)
      .group_by(&:locale)
      .transform_values { |pages| Language::SitemapLandingPageResource.new(pages) }

    lesson_infos_by_locale = Language::Lesson::Version::Info
      .with_locale(ordered_locales)
      .joins(:language)
  .where("language_lesson_version_infos.language_version_id = languages.current_version_id")
      .includes(:lesson, :version)
      .select(
        :locale,
        :name,
        :language_id,
        :language_lesson_id,
        :version_id,
        "language_lesson_versions.natural_order",
        "language_lesson_versions.lesson_id",
        "language_lessons.slug"
      )
      .distinct
      .order("language_lesson_versions.natural_order ASC")
      .group_by(&:locale)
      .transform_values { |infos| infos.group_by(&:language_id) }

    lesson_resources_by_locale_and_language_id = lesson_infos_by_locale.transform_values do |infos_by_language_id|
      infos_by_language_id.transform_values { |infos| Language::SitemapLessonResource.new(infos) }
    end

    blog_post_resources_by_locale = BlogPost.published
      .select(:id, :slug, :name, :locale)
      .order(id: :desc)
      .group_by(&:locale)
      .transform_values { |posts| SitemapBlogPostResource.new(posts) }

      language_category_resources_by_locale = Language::Category
      .select(:id, :name, :slug, :locale)
      .group_by(&:locale)
      .transform_values { |categories| Language::SitemapCategoryResource.new(categories) }

    title = t(".title")

    seo_tags = {
      title:
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      title:,
      orderedLocales: ordered_locales,
      landingPagesByLocale: language_landing_page_resources_by_locale,
      lessonsByLocaleAndLanguageId: lesson_resources_by_locale_and_language_id,
      blogPostsByLocale: blog_post_resources_by_locale,
      categoriesByLocale: language_category_resources_by_locale
    }
  end
end
