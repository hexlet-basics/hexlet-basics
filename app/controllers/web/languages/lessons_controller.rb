# typed: strict

class Web::Languages::LessonsController < Web::Languages::ApplicationController
  allow_unauthenticated_access

  # before_action :authenticate_user!, only: [ :next_lesson ]
  # before_action :redirect_archived_language

  sig { returns(T.untyped) }
  def show
    lesson = resource_language.lessons.find_by(slug: params[:id])

    unless lesson
      f(:lesson_not_found, type: :info)
      redirect_to language_path(resource_language.slug)
      return
    end

    lesson_version = resource_language.current_lesson_versions.find_by!(lesson: lesson)
    lesson_info = lesson_version.infos.includes(language: :current_version).find_by(locale: I18n.locale)

    unless lesson_info
      f(:lesson_not_found, type: :info)
      redirect_to language_path(resource_language.slug)
      return
    end
    # language_info = resource_language.current_version.infos.find_by!(locale: I18n.locale)

    next_lesson = lesson_version.next_lesson
    next_lesson_info = next_lesson ? next_lesson.infos.includes(language: :current_version).find_by!(locale: I18n.locale) : nil
    prev_lesson_version = lesson_version.prev_lesson
    prev_lesson_info = prev_lesson_version ? prev_lesson_version.infos.includes(language: :current_version).find_by!(locale: I18n.locale) : nil

    user = current_user
    # Dynamic creation, because user can start from any lesson directly
    lesson_member =
      if user
        locale = resource_language_landing_page.locale
        result = CourseProgressService.start_lesson(user:, language: resource_language, lesson:, locale:)

        case result
        when Typed::Success
          payload = result.payload
          js_events(payload.events)
          payload.lesson_member
        when Typed::Failure
          nil
        end
      end

    ai_chat = lesson_member && AiChat.find_or_create_by!(
      user: lesson_member.user,
      language_lesson_member: lesson_member
    )

    title = t(
      ".title",
      lesson_name: lesson_info,
      language_name: resource_language_landing_page.name
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
      .includes(:language, :lesson, :version)
      .order(language_lesson_versions: { natural_order: :asc })

    can_create_assistant_message = AiMessagePolicy.new(
      current_user,
      AiMessage
    ).create?

    previous_messages = ai_chat && ai_chat.ai_messages.where(role: %w[user assistant]).order(:id)

    render inertia: true, props: {
      aiChat: ai_chat && AiChatResource.new(ai_chat),
      previousMessages: previous_messages && AiMessageResource.new(previous_messages),
      canCreateAssistantMessage: can_create_assistant_message,
      course: LanguageResource.new(resource_language),
      landingPage: Language::LandingPageForListsResource.new(resource_language_landing_page),
      courseCategory: resource_language_landing_page.language_categories.any? && Language::CategoryResource.new(resource_language_landing_page.language_categories.first),
      lesson: Language::LessonResource.new(lesson_info),
      nextLesson: next_lesson_info && Language::LessonResource.new(next_lesson_info),
      prevLesson: prev_lesson_info && Language::LessonResource.new(prev_lesson_info),
      lessons: Language::LessonResource.new(lessons_infos),
      lessonMember: lesson_member && Language::Lesson::MemberResource.new(lesson_member)
    }
  end
end
