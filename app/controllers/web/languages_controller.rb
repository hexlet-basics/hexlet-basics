class Web::LanguagesController < Web::ApplicationController
  before_action :authenticate_user!, only: [ :success ]

  def show
    language = T.must(landing_page.language)
    # language_info = language.current_version.infos.find_by!(locale: I18n.locale)
    #
    # @builder = CourseSchema.to_builder(@language, @language_version_info)

    language_modules_infos = language.current_module_infos.with_locale
      .includes(version: :module).joins(:version).order(order: :asc)

    language_lessons_infos = language.current_lesson_infos.with_locale
      .includes(:lesson, version: [ :module, :module_version, :lesson ])
      .order(language_lesson_versions: { natural_order: :asc })
    language_lessons_info_resources = language_lessons_infos.map { |info| Language::LessonResource.new(info) }
    lesson_resources_by_module_id = language_lessons_info_resources.group_by do |resource|
      resource.object.version.module.id
    end

    language_member = nil
    next_lesson_info = nil
    if !current_user.guest?
      language_member = language.members.find_by(user: current_user)
      next_lesson_info = language_member&.next_lesson_info
    end

    first_lesson_info = language.current_lesson_infos
      .joins(:lesson).merge(Language::Lesson.ordered).first!

    # recommended_language_pages = Language::LandingPage.with_locale
    #   .where(main: true).where(listed: true)
    #   .includes({ language: [ :current_version, { cover_attachment: :blob } ] })
    #   .order("RANDOM()").excluding(landing_page).limit(4)

    image_url = language.cover.attached? && view_context.rails_representation_url(language.cover.variant(:list))
    seo_tags = {
      title: landing_page.meta_title,
      # keywords: language_info.keywords.join(", "),
      description: landing_page.meta_description,
      canonical: language_url(language.slug),
      image_src: image_url,
      og: {
        title: landing_page.meta_title,
        type: "website",
        description: landing_page.meta_description,
        url: language_url(landing_page.slug),
        image: image_url,
        locale: I18n.locale
      }
    }
    set_meta_tags seo_tags

    # @switching_locales.each do |locale,|
    #   if @language.current_version.infos.exists?(locale: locale)
    #     @switching_locales[locale] = language_path(@language.slug, locale: AppHost.locale_for_url(locale))
    #   end
    # end

    # reviews = Review.published.with_locale
    #   .includes([ :user, :language ])
    #   .where(language: landing_page.language)
    #   .order(
    #     Review.arel_table[:pinned].desc.nulls_last,
    #     id: :desc
    #   ).limit(6)

    render inertia: true, props: {
      courseLandingPage: Language::LandingPageResource.new(landing_page),
      courseLandingPageQnaItems: Language::LandingPageQnaItemResource.new(landing_page.qna_items),
      course: LanguageResource.new(language),
      courseCategory: Language::CategoryResource.new(landing_page.language_category),
      firstLesson: Language::LessonResource.new(first_lesson_info),
      nextLesson: next_lesson_info && Language::LessonResource.new(next_lesson_info),
      courseModules: Language::ModuleResource.new(language_modules_infos),
      lessonsByModuleId: lesson_resources_by_module_id,
      courseMember: language_member && Language::MemberResource.new(language_member)
      # reviews: ReviewResource.new(reviews)
    }
  end

  def success
    language_member = landing_page.language.members.find_by(user: current_user)
    unless language_member.finished?
      f("error")
      redirect_to language_path(landing_page.slug)
      return
    end

    render inertia: true, props: {
      courseLandingPage: Language::LandingPageResource.new(landing_page),
      lead: LeadForm.new
    }
  end

  private

  def landing_page
    @language_page ||= Language::LandingPage.published.find_by!(locale: I18n.locale, slug: params[:id])
  end
end
