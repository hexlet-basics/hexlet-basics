# frozen_string_literal: true

class Web::Languages::LessonsController < Web::Languages::ApplicationController
  def show
    @lesson = resource_language.lessons.find_by!(slug: params[:id])
    @description = @lesson.infos.find_by!(locale: I18n.locale)

    gon.language = resource_language.to_s
    gon.locale = I18n.locale
    gon.lesson = @lesson.current_version
  end
end
