# frozen_string_literal: true

class Api::Languages::LessonsController < Api::Languages::ApplicationController
  def index
    @lesson_infos = resource_language.current_lesson_infos
                                     .joins(version: :lesson)
                                     .with_locale
                                     .order("language_lesson_versions.natural_order")
  end

  def show
    @lesson = resource_language_version.lessons.find(params[:id])
    @lesson_version = resource_language_version.lesson_versions.find_by!(lesson: @lesson)
    @lesson_info = @lesson_version.infos.find_by!(locale: I18n.locale)
  end
end
