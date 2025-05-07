module SurveyConcern
  extend ActiveSupport::Concern

  def run_survey_if_needed
    return if current_user.guest?

    if request.get? && I18n.locale == :ru
      member = current_user.survey_scenario_members.started.first
      return unless member

      next_survey = member.next_survey
      if next_survey
        target_path = view_context.scenario_survey_path(member.scenario, next_survey)
        if request.path != target_path
          redirect_to view_context.scenario_survey_path(member.scenario, next_survey, from: request.path)
        end
      else
        member.state = :finished
        member.save!
      end
    end
  end
end
