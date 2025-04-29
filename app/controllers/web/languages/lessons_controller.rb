# frozen_string_literal: true

class Web::Languages::LessonsController < Web::Languages::ApplicationController
  # before_action :authenticate_user!, only: [ :next_lesson ]

  def show
    lesson = resource_language.lessons.find_by(slug: params[:id])

    unless lesson
      f(:lesson_not_found, type: :info)
      redirect_to language_path(resource_language.slug)
      return
    end

    lesson_version = resource_language.current_lesson_versions.find_by!(lesson: lesson)
    lesson_info = lesson_version.infos.find_by!(locale: I18n.locale)
    # language_info = resource_language.current_version.infos.find_by!(locale: I18n.locale)

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
          locale: resource_language_landing_page.locale.to_sym
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
          locale: resource_language_landing_page.locale.to_sym
        }
        event = LessonStartedEvent.new(data: event_data)

        publish_event(event, current_user)
        event_to_js(event)
      end
    end

    title = t(
      ".title",
      lesson: lesson_info,
      language: resource_language_landing_page.header
    ).squish
    description = view_context.truncate("[#{resource_language.current_version}] — #{lesson_info} — #{lesson_info.theory}", length: 220)

    image_url = resource_language.cover.attached? && view_context.rails_representation_url(resource_language.cover.variant(:list))
    seo_tags = {
      title:,
      description:,
      canonical: language_lesson_url(resource_language.slug, lesson.slug),
      # amphtml: language_lesson_url(@lesson.language.slug, @lesson.slug, format: "amp", only_path: false),
      image_src: image_url,
      og: {
        type: "article",
        locale: I18n.locale,
        title: title,
        url: language_lesson_url(resource_language.slug, lesson.slug),
        image: image_url
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

    can_create_assistant_message = Language::Lesson::Member::MessagePolicy.new(
      current_user,
      Language::Lesson::Member::Message
    ).create?

    # raise resource_language.category.inspect
    render inertia: true, props: {
      previousMessages: Language::Lesson::Member::MessageResource.new(lesson_member.messages.order(id: :asc)),
      canCreateAssistantMessage: can_create_assistant_message,
      course: LanguageResource.new(resource_language),
      landingPage: Language::LandingPageForListsResource.new(resource_language_landing_page),
      courseCategory: Language::CategoryResource.new(resource_language_landing_page.language_category),
      lesson: Language::LessonResource.new(lesson_info),
      nextLesson: next_lesson_info && Language::LessonResource.new(next_lesson_info),
      prevLesson: prev_lesson_info && Language::LessonResource.new(prev_lesson_info),
      lessons: Language::LessonResource.new(lessons_infos),
      lessonMember: lesson_member && Language::Lesson::MemberResource.new(lesson_member)
    }
  end
end
