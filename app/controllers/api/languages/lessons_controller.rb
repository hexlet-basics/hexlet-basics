# frozen_string_literal: true

class Api::Languages::LessonsController < Api::Languages::ApplicationController
  def show
    @lesson = resource_language_version.lessons.find(params[:id])
    @lesson_version = resource_language_version.lesson_versions.find_by!(lesson: @lesson)
    @lesson_info = @lesson_version.infos.find_by!(locale: I18n.locale)

    respond_with @lesson_version
  end
end
