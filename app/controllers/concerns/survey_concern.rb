module SurveyConcern
  extend ActiveSupport::Concern

  # included do
  # end

  def run_survey_if_needed
    if request.get? && I18n.locale == :ru
      answer = Survey::Answer.requested.order(id: :asc).find_by user: current_user
      return unless answer

      target_path = survey_path(answer.survey.slug)
      unless request.path == target_path
        redirect_to survey_path(answer.survey.slug, from: request.path)
      end
    end
  end
end
