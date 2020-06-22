# frozen_string_literal: true

class Web::Languages::Modules::LessonsController < Web::Languages::Modules::ApplicationController
  def show
    @lesson = resource_module.lessons.find_by!(slug: params[:id])
    @description = @lesson.descriptions.find_by!(locale: I18n.locale)

    current_language_version = @lesson.language.current_version
    current_lesson_version = @lesson.current_version
    gon.language = current_language_version.name
    gon.locale = I18n.locale
    gon.lesson = current_lesson_version

    render :show, layout: 'lesson'
  end
end
