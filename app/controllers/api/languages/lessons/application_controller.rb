# frozen_string_literal: true

class Api::Languages::Lessons::ApplicationController < Api::Languages::ApplicationController
  helper_method :resource_lesson

  def resource_lesson
    @resource_lesson ||= Language::Lesson.find(params[:lesson_id])
  end
end
