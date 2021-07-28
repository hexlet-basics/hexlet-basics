# frozen_string_literal: true

class Web::Languages::LessonsController < Web::Languages::ApplicationController
  before_action :authenticate_user!, only: [:next_lesson]

  def show
    # raise params[:controller].inspect
    @lesson = resource_language.lessons.find_by(slug: params[:id])
    unless @lesson
      f(:lesson_not_found, type: :info)
      redirect_to language_path(resource_language.slug)
      return
    end
    @lesson_version = resource_language.current_lesson_versions.find_by!(lesson: @lesson)
    @info = @lesson_version.infos.find_by!(locale: I18n.locale)
    @language_lessons_count = resource_language.current_lessons.count

    title t("human_languages.#{@resource_language}")
    title @info

    if current_user.guest?
      gon.lesson_member = Language::Lesson::MemberFake.new
    else
      lesson_member = @lesson.members.find_or_create_by!(language: resource_language, user: current_user)
      gon.lesson_member = lesson_member
    end

    gon.language = resource_language.to_s
    gon.lesson_version = @lesson_version
    gon.lesson = @lesson

    set_meta_tags canonical: language_lesson_url(@lesson.language.slug, @lesson.slug),
                  amphtml: language_lesson_url(@lesson.language.slug, @lesson.slug, format: 'amp', only_path: false),
                  title: [@info, t("human_languages.#{@resource_language}")].join(' | ')
  end

  def next_lesson
    language_slug = params[:language_id]
    lesson = resource_language.lessons.find_by!(slug: params[:id])
    lesson_version = resource_language.current_lesson_versions.find_by!(lesson: lesson)

    language_member = current_user.language_members.find_by! language: resource_language

    next_lesson = lesson_version.next_lesson

    # TODO Добавить сериализацию language, lesson, language_member
    js_event_options = {
      user: current_user,
      language: resource_language.to_hash,
      lesson: next_lesson&.to_hash,
      language_member: language_member.to_hash,
      lessons_started: current_user.lesson_members.where(language: resource_language).count,
      lessons_finished: current_user.lesson_members.where(language: resource_language).finished.count
    }
    js_event :next_lesson, js_event_options

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

    prev_lesson = lesson_version.prev_lesson

    if prev_lesson.nil?
      redirect_to language_path(language_slug)
    else
      redirect_to language_lesson_path(language_slug, prev_lesson.slug)
    end
  end
end
