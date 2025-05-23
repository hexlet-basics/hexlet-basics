class Web::Admin::SurveysController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "created_at", "so" => "0")
    search = Survey.where(locale: I18n.locale).ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      surveys: SurveyCrudResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  def new
    survey = Admin::SurveyForm.new

    render inertia: true, props: {
      surveyItems: [],
      surveyDto: SurveyCrudResource.new(survey)
    }
  end

  def edit
    survey = Admin::SurveyForm.find(params[:id])

    render inertia: true, props: {
      surveyDto: SurveyCrudResource.new(survey),
      surveyItems: Survey::ItemCrudResource.new(survey.items)
    }
  end

  def create
    survey = Admin::SurveyForm.new(params[:survey])
    survey.locale = I18n.locale

    if survey.save
      f(:success)
      redirect_to edit_admin_survey_url(survey)
    else
      f(:error)
      redirect_to_inertia new_admin_survey_url, survey
    end
  end

  def update
    survey = Admin::SurveyForm.find(params[:id])
    survey.locale = I18n.locale

    if survey.update(params[:survey])
      f(:success)
    else
      f(:error)
    end

    redirect_to_inertia edit_admin_survey_url(survey), survey
  end
end
