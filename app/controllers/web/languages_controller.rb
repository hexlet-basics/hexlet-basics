# frozen_string_literal: true

class Web::LanguagesController < Web::ApplicationController
  def show
    @language = Language.find_by!(slug: params[:id])
    @current_modules = @language.current_modules
                                .includes(:module)
                                .order(:order)
                                .eager_load(:lesson_versions)
                                .merge(
                                  Language::Lesson::Version.includes(:lesson).order(:order)
                                )
    @infos_by_module = @language.current_module_infos.with_locale.index_by(&:version_id)
    @infos_by_lesson = @language.current_lesson_infos.with_locale.index_by(&:version_id)
    finished_members = if current_user.guest?
                         []
                       else
                         current_user.lesson_members.finished
                       end
    @finished_members_by_lesson = finished_members.index_by(&:lesson_id)

    @first_lesson = @language.current_lessons.order(:natural_order).limit(1).first
    @next_lesson = if current_user.guest?
                     @first_lesson
                   else
                     not_finished_lessons = @language.current_lessons.where.not(lesson_id: finished_members.map(&:lesson_id)).order(:natural_order).limit(1)
                     not_finished_lessons.first
                   end
  end
end
