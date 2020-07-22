# frozen_string_literal: true

class Web::Languages::LessonsController < Web::Languages::ApplicationController
  def show
    @lesson = resource_language.lessons.find_by!(slug: params[:id])
    lesson_version = resource_language.current_lesson_versions.find_by!(lesson: @lesson)
    @info = @lesson.infos.find_by!(locale: I18n.locale)
    @lesson.members.find_or_create_by!(language: @lesson.language, user: current_user) unless current_user.guest?

    gon.language = resource_language.to_s
    gon.locale = I18n.locale
    gon.lesson = lesson_version
  end
end
