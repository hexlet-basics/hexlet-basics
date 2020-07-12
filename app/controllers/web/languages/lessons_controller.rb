# frozen_string_literal: true

class Web::Languages::LessonsController < Web::Languages::ApplicationController
  def show
    @lesson = resource_language.lessons.find_by!(slug: params[:id])
    @info = @lesson.infos.find_by!(locale: I18n.locale)
    lesson_member = Lesson::MemberService.ensure_member!(@lesson, current_user)

    gon.language = resource_language.to_s
    gon.locale = I18n.locale
    gon.lesson = lesson_member.lesson_version
  end
end
