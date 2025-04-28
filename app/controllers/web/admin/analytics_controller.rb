class Web::Admin::AnalyticsController < Web::Admin::ApplicationController
  def index
  end

  def surveys
    q = ransack_params("sf" => "created_at", "so" => "0")
    search = UserSurveyPivot.ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      userSurveyAnswers: UserSurveyPivotResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }
  end
end
