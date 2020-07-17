# frozen_string_literal: true

class Web::Languages::LessonsController < Web::Languages::ApplicationController
  def show
    @lesson = resource_language.lessons.find_by!(slug: params[:id])
    @info = @lesson.infos.find_by!(locale: I18n.locale)
    lesson_version = if current_user.guest?
                       nil
                     else
                       member = LessonMemberMutator.upsert_member!(@lesson, current_user)
                       member.lesson_version
                     end

    gon.language = resource_language.to_s
    gon.locale = I18n.locale
    gon.lesson = lesson_version
  end
end
