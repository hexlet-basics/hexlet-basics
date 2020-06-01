# frozen_string_literal: true

class Web::Languages::Modules::LessonsController < Web::Languages::Modules::ApplicationController
  def show
    @lesson = Language::Module::Lesson.find_by(id: params[:id])
    @description = @lesson.descriptions.where(locale: I18n.locale).first
    render :show, layout: 'lesson'
  end
end
