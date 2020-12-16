# frozen_string_literal: true

class Web::Languages::LessonsController < Web::Languages::ApplicationController
  def show
    @lesson = resource_language.lessons.find_by!(slug: params[:id])
    @lesson_version = resource_language.current_lesson_versions.find_by!(lesson: @lesson)
    @info = @lesson_version.infos.find_by!(locale: I18n.locale)
    @language_lessons_count = resource_language.current_lessons.count

    return if current_user.guest?

    lesson_member = @lesson.members.find_or_create_by!(language: resource_language, user: current_user)

    gon.lesson_member = lesson_member
    gon.language = resource_language.to_s
    gon.lesson_version = @lesson_version
    gon.lesson = @lesson

    title @resource_language
    title @info
  end

  def next_lesson
    language_slug = params[:language_id]
    lesson = resource_language.lessons.find_by!(slug: params[:id])
    lesson_version = resource_language.current_lesson_versions.find_by!(lesson: lesson)

    next_lesson = lesson_version.next_lesson

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
