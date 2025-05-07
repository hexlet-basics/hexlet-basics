class Web::Admin::SurveyScenariosController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "created_at", "so" => "0")
    search = Survey::Scenario.where(locale: I18n.locale).ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      surveyScenarios: Survey::ScenarioCrudResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end

  def new
    survey_scenario = Admin::SurveyScenarioForm.new
    surveys = Survey.where(locale: I18n.locale)

    render inertia: true, props: {
      surveysItems: Survey::ItemCrudResource.new(Survey::Item.order(order: :asc)),
      surveyScenarioDto: Survey::ScenarioCrudResource.new(survey_scenario),
      surveys: SurveyResource.new(surveys)
    }
  end

  def edit
    survey_scenario = Admin::SurveyScenarioForm.find(params[:id])
    surveys = Survey.where(locale: I18n.locale)

    render inertia: true, props: {
      surveyScenarioDto: Survey::ScenarioCrudResource.new(survey_scenario),
      surveysItems: Survey::ItemCrudResource.new(Survey::Item.includes(:survey).order(order: :asc)),
      surveys: SurveyResource.new(surveys)
    }
  end

  def create
    survey_scenario = Admin::SurveyScenarioForm.new(params[:survey_scenario])
    survey_scenario.locale = I18n.locale

    if survey_scenario.save
      f(:success)
      redirect_to view_context.edit_admin_survey_scenario_url(survey_scenario)
    else
      f(:error)
      redirect_to_inertia view_context.new_admin_survey_scenario_url, survey_scenario
    end
  end

  def update
    survey_scenario = Admin::SurveyScenarioForm.find(params[:id])
    survey_scenario.locale = I18n.locale

    if survey_scenario.update(params[:survey_scenario])
      f(:success)
    else
      f(:error)
    end

    redirect_to_inertia view_context.edit_admin_survey_scenario_url(survey_scenario), survey_scenario
  end
end
