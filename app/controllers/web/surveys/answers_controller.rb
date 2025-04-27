class Web::Surveys::AnswersController < Web::ApplicationController
  before_action :authenticate_user!

  def create
    survey = Survey.find_by! slug: params[:survey_id]
    answer = survey.answers.find_by! user: current_user
    answer.survey_item_id = params[:survey_answer][:survey_item_id]
    answer.state = "fulfilled"
    answer.save!

    requested_answer = current_user.survey_answers.requested.first

    if requested_answer
      redirect_to survey_path(requested_answer.survey.slug, from: params[:from])
      return
    end

    url = params[:from].presence || root_path
    redirect_to url
  end
end
