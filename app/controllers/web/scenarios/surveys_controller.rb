class Web::Scenarios::SurveysController < Web::Scenarios::ApplicationController
  def show
    survey = resource_scenario.surveys.where(locale: I18n.locale).find params[:id]
    member = resource_scenario.members.find_by user: current_user
    unless member
      redirect_to view_context.root_path
      return
    end

    answer = survey.answers.find_by user: current_user
    if answer
      f(:success)
      redirect_to params[:from].presence || view_context.root_path
      return
    end

    survey_answers = survey.items.active.order(order: :asc)

    render inertia: true, props: {
      survey: SurveyResource.new(survey),
      scenario: Survey::ScenarioResource.new(resource_scenario),
      surveyItems: Survey::ItemResource.new(survey_answers)
    }
  end
end
