# frozen_string_literal: true

class Web::LanguagesController < Web::ApplicationController
  def show
    @language = Language.find_by!(slug: params[:id])
    @current_module_versions = @language.current_module_versions
                                        .includes(:module)
                                        .order(:order)
                                        .eager_load(:lesson_versions)
                                        .merge(
                                          Language::Lesson::Version.includes(:lesson).order(:order)
                                        )
    @infos_by_module = @language.current_module_infos.with_locale.index_by(&:version_id)
    @infos_by_lesson = @language.current_lesson_infos.with_locale.index_by(&:version_id)

    finished_members = current_user.finished_members_for_language(@language)
    @finished_members_by_lesson = finished_members.index_by(&:lesson_id)

    @first_lesson = @language.current_lesson_versions.order(:natural_order).first.lesson
    @next_lesson = @language.next_lesson_for_user(current_user)
  end
end
