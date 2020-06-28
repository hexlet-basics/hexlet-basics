# frozen_string_literal: true

class Web::Languages::Modules::LessonsController < Web::Languages::Modules::ApplicationController
  def show
    @lesson = resource_module.lessons.find_by!(slug: params[:id])
    @description = @lesson.descriptions.find_by!(locale: I18n.locale)
    lesson_member = LessonMember::EnsureService.new(@lesson, current_user).execute

    gon.language = resource_language.to_s
    gon.locale = I18n.locale
    gon.lesson = lesson_member.lesson_version

    render :show, layout: 'lesson'
  end
end
