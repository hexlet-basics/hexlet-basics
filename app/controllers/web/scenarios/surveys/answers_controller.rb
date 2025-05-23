class Web::Scenarios::Surveys::AnswersController < Web::Scenarios::ApplicationController
  def create
    survey = resource_scenario.surveys.where(locale: I18n.locale).find params[:survey_id]
    member = resource_scenario.members.find_by user: current_user
    unless member
      redirect_to view_context.root_path
      return
    end

    answer = survey.answers.find_or_initialize_by user: current_user
    answer.survey_item_id = params[:survey_answer][:survey_item_id]
    answer.save!

    current_user.tag_list.add(answer.survey_item.tag_list)
    current_user.save!

    next_survey = member.next_survey

    event_data = {
      survey_scenario_member_id: member.id,
      survey_answer_id: answer.id,
      next_survey_id: next_survey&.id
    }
    event = SurveyAnsweredEvent.new(data: event_data)
    publish_event(event, current_user)

    if next_survey
      redirect_to view_context.scenario_survey_path(resource_scenario, next_survey, from: params[:from])
      return
    end

    url = params[:from].presence || root_path
    redirect_to url
  end
end
