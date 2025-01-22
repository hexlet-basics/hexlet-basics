class Web::LanguagesController < Web::ApplicationController
  def show
    language = Language.find_by!(slug: params[:id])
    language_info = language.current_version.infos.find_by!(locale: I18n.locale)
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
    if !current_user.guest?
      language_member = language.members.find_by(user: current_user)
    end

    first_lesson_info = language.current_lesson_infos
      .joins(:lesson).merge(Language::Lesson.ordered).first

    finished_lesson_ids = current_user.finished_lessons_for_language(language).pluck(:id)

    next_lesson_info = nil
    if !current_user.guest?
      next_lesson_info = language.current_lesson_infos.where.not(language_lesson_id: finished_lesson_ids)
        .joins(:lesson).merge(Language::Lesson.ordered).first
    end

    recommendedCourses = Language::Version::Info.current.with_locale.order("RANDOM()").excluding(language_info).limit(4)

    # @blog_posts = @language.blog_posts.published
    #
    # gon.language = @language.slug
    #
    # human_language_header = [ @language.current_version.name, @language.learn_as.text ].join(' ')
    # @header = @language_version_info.header || human_language_header
    # title = @language_version_info.title || @header
    # description = @language_version_info.seo_description || @language_version_info.description

    seo_tags = {
      title: language_info.title,
      keywords: language_info.keywords.join(", "),
      description: language_info.seo_description,
      canonical: language_url(language.slug),
      image_src: view_context.vite_asset_url("images/#{language.slug}.png"),
      og: {
        title: language_info.title,
        type: "website",
        description: language_info.seo_description,
        url: language_url(language.slug),
        image: view_context.vite_asset_url("images/#{language.slug}.png"),
        locale: I18n.locale
      }
    }
    set_meta_tags seo_tags

    # @switching_locales.each do |locale,|
    #   if @language.current_version.infos.exists?(locale: locale)
    #     @switching_locales[locale] = language_path(@language.slug, locale: AppHost.locale_for_url(locale))
    #   end
    # end

    render inertia: true, props: {
      course: LanguageResource.new(language_info),
      finishedLessonIds: finished_lesson_ids,
      courseCategory: Language::CategoryResource.new(language.category),
      firstLesson: Language::LessonResource.new(first_lesson_info),
      nextLesson: !current_user.guest? && Language::LessonResource.new(next_lesson_info),
      courseModules: Language::ModuleResource.new(language_modules_infos),
      lessonsByModuleId: lesson_resources_by_module_id,
      courseMember: !current_user.guest? && Language::MemberResource.new(language_member),
      recommendedCourses: LanguageResource.new(recommendedCourses)
    }
  end
end
