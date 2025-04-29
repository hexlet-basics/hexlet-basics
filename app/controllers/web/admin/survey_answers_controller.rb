class Web::Admin::SurveyAnswersController < Web::Admin::ApplicationController
  def index
    q = ransack_params("sf" => "created_at", "so" => "0")
    search = Survey::Answer.ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      surveyAnswers: Survey::AnswerResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end
end
