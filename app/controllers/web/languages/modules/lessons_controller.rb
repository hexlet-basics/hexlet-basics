# frozen_string_literal: true

class Web::Languages::Modules::LessonsController < Web::Languages::Modules::ApplicationController
  def show
    @lesson = resource_module.lessons.find_by!(slug: params[:id])
    @description = @lesson.descriptions.find_by!(locale: I18n.locale)

    render :show, layout: 'lesson'
  end
end
