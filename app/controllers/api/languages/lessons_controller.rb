# frozen_string_literal: true

class Api::Languages::LessonsController < Api::Languages::ApplicationController
  def index
    @lessons = resource_language_version.lessons

    @current_module_versions = resource_language.current_module_versions
                                                .eager_load(:lesson_versions)
                                                .joins(:infos)
                                                .merge(Language::Module::Version::Info.with_locale)
                                                .merge(Language::Lesson::Version.includes(:lesson).order(:order))
                                                .order(:order)

    @infos_by_lesson = resource_language.current_lesson_infos
                                        .includes(version: :lesson)
                                        .with_locale
                                        .index_by(&:version_id)

    respond_with @lessons
  end

  def show
    @lesson = resource_language_version.lessons.find(params[:id])
    @lesson_version = resource_language_version.lesson_versions.find_by!(lesson: @lesson)
    @lesson_info = @lesson_version.infos.find_by!(locale: I18n.locale)

    respond_with @lesson_version
  end
end
