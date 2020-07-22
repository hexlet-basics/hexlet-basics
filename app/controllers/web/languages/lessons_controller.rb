# frozen_string_literal: true

class Web::Languages::LessonsController < Web::Languages::ApplicationController
  def show
    @lesson = resource_language.lessons.find_by!(slug: params[:id])
    lesson_version = resource_language.current_lesson_versions.find_by!(lesson: @lesson)
    @info = @lesson.infos.find_by!(locale: I18n.locale)

    return if current_user.guest?

    lesson_member = @lesson.members.find_or_create_by!(language: @lesson.language, user: current_user)

    gon.next_lesson = @lesson.next_lesson
    gon.prev_lesson = @lesson.prev_lesson
    gon.lesson_member = lesson_member
    gon.language = resource_language.to_s
    gon.locale = I18n.locale
    gon.lesson_version = lesson_version
  end
end
