# frozen_string_literal: true

class Web::LanguagesController < Web::ApplicationController
  def show
    language = Language.find_by!(slug: params[:id])
    language_info = language.current_version.infos.find_by!(locale: I18n.locale)
    #
    # @builder = CourseSchema.to_builder(@language, @language_version_info)
    #
    # if @language.progress_in_development?
    #   f('.language_in_development_html', type: :info, values: { language: @language.to_s, link_to_repo: ExternalLinks.source_code_curl, link_to_recommendations: page_path(:authors) }, now: true)
    # end
    #
    # language_module_versions = language.current_module_versions
    #                                     .includes(:module)
    #                                     .order(:order)
    #                                     .eager_load(:lesson_versions)
    #                                     .joins(:infos)
    #                                     .merge(Language::Module::Version::Info.with_locale)
    #                                     .merge(Language::Lesson::Version.includes(:lesson).order(:order))

    language_modules_infos = language.current_module_infos.with_locale
    # .merge(Language::Module::Version::Info.with_locale)

    # language_module_resources = language_modules_infos.map { |info| Language::ModuleResource.new(info) }
    # language_module_resources_by_id = language_module_resources.index_by { |r| r.object.id }

    #
    # @infos_by_module = @language.current_module_infos.with_locale.index_by(&:version_id)
    language_lessons_infos = language.current_lesson_infos.with_locale
    language_lessons_info_resources = language_lessons_infos.map { |info| Language::LessonResource.new(info) }
    lesson_resources_by_module_id = language_lessons_info_resources.group_by do |resource|
      resource.object.version.module.id
    end

    #
    # @finished_lessons_by_id = current_user.finished_lessons_for_language(@language).index_by(&:id)
    # @language_member = @language.members.find_by(user: current_user) || Language::MemberFake.new
    #

    first_lesson_info = language.current_lesson_infos
      .joins(:lesson).merge(Language::Lesson.ordered).first
    next_lesson = current_user.not_finished_lessons_for_language(language)
      .joins(:lesson).merge(Language::Lesson.ordered).first
    #
    recommendedCourses = Language::Version::Info.with_locale.order("RANDOM()").excluding(language_info).limit(4)
    # @blog_posts = @language.blog_posts.published
    #
    # gon.language = @language.slug
    #
    # human_language_header = [ @language.current_version.name, @language.learn_as.text ].join(' ')
    # @header = @language_version_info.header || human_language_header
    # title = @language_version_info.title || @header
    # description = @language_version_info.seo_description || @language_version_info.description
    #
    # seo_tags = {
    #   title: title,
    #   keywords: @language_version_info.keywords.join(', '),
    #   description: description,
    #   canonical: language_url(@language.slug),
    #   image_src: view_context.asset_url("#{@language.slug}.png"),
    #   og: {
    #     title: title,
    #     type: 'website',
    #     description: description,
    #     url: language_url(@language.slug),
    #     image: view_context.asset_url("#{@language.slug}.png"),
    #     locale: I18n.locale
    #   }
    # }
    # set_meta_tags seo_tags
    #
    # @switching_locales.each do |locale,|
    #   if @language.current_version.infos.exists?(locale: locale)
    #     @switching_locales[locale] = language_path(@language.slug, locale: AppHost.locale_for_url(locale))
    #   end
    # end

    render inertia: true, props: {
      course: LanguageResource.new(language_info),
      courseCategory: Language::CategoryResource.new(language.category),
      firstLesson: Language::LessonResource.new(first_lesson_info),
      nextLesson: Language::LessonResource.new(next_lesson),
      courseModules: Language::ModuleResource.new(language_modules_infos),
      lessonsByModuleId: lesson_resources_by_module_id,
      recommendedCourses: LanguageResource.new(recommendedCourses)
    }
  end
end
