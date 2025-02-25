# frozen_string_literal: true

class Web::Languages::LessonsController < Web::Languages::ApplicationController
  before_action :authenticate_user!, only: [ :next_lesson ]

  def show
    lesson = resource_language.lessons.find_by(slug: params[:id])

    unless lesson
      f(:lesson_not_found, type: :info)
      redirect_to language_path(resource_language.slug)
      return
    end

    lesson_version = resource_language.current_lesson_versions.find_by!(lesson: lesson)
    lesson_info = lesson_version.infos.find_by!(locale: I18n.locale)
    language_info = resource_language.current_version.infos.find_by!(locale: I18n.locale)

    next_lesson_version = lesson_version.next_lesson_version
    next_lesson_info = next_lesson_version ? next_lesson_version.infos.find_by!(locale: I18n.locale) : nil
    prev_lesson_version = lesson_version.prev_lesson_version
    prev_lesson_info = prev_lesson_version ? prev_lesson_version.infos.find_by!(locale: I18n.locale) : nil

    lesson_member = nil

    if !current_user.guest?
      # Dynamic creation, because user can start from any lesson directly
      language_member = resource_language.members.find_or_initialize_by(user: current_user)
      if language_member.new_record?
        language_member.save!
        event_data = {
          slug: resource_language.slug,
          locale: language_info.locale
        }
        event = CourseStartedEvent.new(data: event_data)
        publish_event(event, current_user)
        event_to_js(event)
      end
      lesson_member = language_member.lesson_members.find_or_initialize_by(
        language: resource_language,
        lesson:,
        user: current_user,
      )
      if lesson_member.new_record?
        lesson_member.save!

        event_data = {
          lesson_slug: lesson.slug,
          course_slug: resource_language.slug,
          locale: language_info.locale
        }
        event = LessonStartedEvent.new(data: event_data)

        publish_event(event, current_user)
        event_to_js(event)
      end
    end

    # unless @lesson
    #   f(:lesson_not_found, type: :info)
    #   redirect_to language_path(resource_language.slug)
    #   return
    # end
    #
    # @lesson_version = resource_language.current_lesson_versions.find_by!(lesson: @lesson)
    # @info = @lesson_version.infos.with_locale.sole
    # @language_lessons_count = resource_language.current_lessons.count
    # @lessons_info = resource_language.current_lesson_infos
    #                                  .joins(version: :lesson)
    #                                  .includes(version: :lesson)
    #                                  .with_locale
    #                                  .order('language_lesson_versions.natural_order')
    #
    # if current_user.guest?
    #   gon.lesson_member = Language::Lesson::MemberFake.new
    # else
    #
    #   if language_member.new_record?
    #     language_member.save!
    #     js_event_options = {
    #       user: current_user.serializable_data,
    #       language: resource_language.serializable_data,
    #       language_member: language_member.serializable_data
    #     }
    #     js_event :language_started, js_event_options
    #   end
    #
    #
    #   if lesson_member.new_record?
    #     lesson_member.save!
    #     js_event_options = {
    #       user: current_user.serializable_data,
    #       language: resource_language.serializable_data,
    #       language_member: language_member.serializable_data,
    #       lesson_member: lesson_member.serializable_data,
    #       lesson: @lesson.serializable_data
    #     }
    #     js_event :lesson_started, js_event_options
    #   end
    #
    #   gon.lesson_member = lesson_member
    # end
    #
    # gon.language = resource_language.slug
    # gon.lesson_version = @lesson_version
    # gon.lesson = @lesson
    #
    title = t(".title", lesson: lesson_info, language: language_info).squish
    description = view_context.truncate("[#{resource_language.current_version}] — #{lesson_info} — #{lesson_info.theory}", length: 220)

    seo_tags = {
      title:,
      description:,
      canonical: language_lesson_url(resource_language.slug, lesson.slug),
      # amphtml: language_lesson_url(@lesson.language.slug, @lesson.slug, format: "amp", only_path: false),
      image_src: view_context.vite_asset_url("images/#{resource_language.slug}.png"),
      og: {
        type: "article",
        locale: I18n.locale,
        title: title,
        url: language_lesson_url(resource_language.slug, lesson.slug),
        image: view_context.vite_asset_url("images/#{resource_language.slug}.png")
      }
    }
    set_meta_tags seo_tags

    # @switching_locales.each do |locale,|
    #   if @lesson_version.infos.exists?(locale: locale)
    #     @switching_locales[locale] = language_lesson_url(resource_language.slug, @lesson.slug, locale: AppHost.locale_for_url(locale))
    #   end
    # end

    lessons_infos = resource_language.current_lesson_infos.with_locale
      .includes(:lesson, :version)
      .order(language_lesson_versions: { natural_order: :asc })

    render inertia: true, props: {
      course: LanguageResource.new(language_info),
      courseCategory: Language::CategoryResource.new(resource_language.category),
      lesson: Language::LessonResource.new(lesson_info),
      nextLesson: next_lesson_info && Language::LessonResource.new(next_lesson_info),
      prevLesson: prev_lesson_info && Language::LessonResource.new(prev_lesson_info),
      lessons: Language::LessonResource.new(lessons_infos),
      lessonMember: lesson_member && Language::Lesson::MemberResource.new(lesson_member)
    }
  end

  def next_lesson
    language_slug = params[:language_id]
    lesson = resource_language.lessons.find_by!(slug: params[:id])
    lesson_version = resource_language.current_lesson_versions.find_by!(lesson: lesson)

    # language_member = current_user.language_members.find_by! language: resource_language

    next_lesson = lesson_version.next_lesson_version

    # TODO Добавить сериализацию language, lesson, language_member
    # NOTE Временно отключил и заменил на language_finished
    # js_event_options = {
    #   user: current_user,
    #   language: resource_language.serializable_data,
    #   lesson: next_lesson&.serializable_data,
    #   language_member: language_member.serializable_data,
    #   lessons_started: current_user.lesson_members.where(language: resource_language).count,
    #   lessons_finished: current_user.lesson_members.where(language: resource_language).finished.count
    # }
    # js_event :next_lesson, js_event_options

    if next_lesson.nil?
      redirect_to language_path(language_slug)
    else
      redirect_to language_lesson_path(language_slug, next_lesson.slug)
    end
  end

  def prev_lesson
    language_slug = params[:language_id]
    lesson = resource_language.lessons.find_by!(slug: params[:id])
    lesson_version = resource_language.current_lesson_versions.find_by!(lesson: lesson)

    prev_lesson = lesson_version.prev_lesson_version

    if prev_lesson.nil?
      redirect_to language_path(language_slug)
    else
      redirect_to language_lesson_path(language_slug, prev_lesson.slug)
    end
  end
end
