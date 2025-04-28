class Web::SurveysController < Web::ApplicationController
  before_action :authenticate_user!

  def show
    survey = Survey.where(locale: I18n.locale).find_by! slug: params[:id]

    answer = survey.answers.find_by user: current_user
    if answer.fulfilled?
      f(:success)
      redirect_to params[:from].presence || root_path
      return
    end

    render inertia: true, props: {
      survey: SurveyResource.new(survey),
      surveyItems: Survey::ItemResource.new(survey.items.active.order(order: :asc))
    }
  end
end
