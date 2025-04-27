class Web::SurveysController < Web::ApplicationController
  before_action :authenticate_user!

  def show
    survey = Survey.where(locale: I18n.locale).find_by! slug: params[:id]

    render inertia: true, props: {
      survey: SurveyResource.new(survey),
      surveyItems: Survey::ItemResource.new(survey.items.active.order(order: :asc))
    }
  end
end
