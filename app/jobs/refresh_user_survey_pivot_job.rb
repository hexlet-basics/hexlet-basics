# frozen_string_literal: true

class RefreshUserSurveyPivotJob < ApplicationJob
  def perform
    system("bin/rails analytics:refresh_user_survey_pivot")
  end
end
