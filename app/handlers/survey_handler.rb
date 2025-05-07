class SurveyHandler
  def call(event)
    # locale = event.metadata[:locale]
    return if I18n.locale != :ru
    user_id = event.metadata[:user_id]
    return unless user_id
    user = User.find(user_id)

    case event
    when SurveyAnsweredEvent
      member_id = event.data.fetch(:survey_scenario_member_id, 0)
      member = user.survey_scenario_members.find member_id
      unless member.next_survey
        member.state = :finished
        member.save!
      end
    else
      count = event.data.fetch(:occurrence_count, 1)
      survey_item_arel = Survey::Scenario.arel_table[:survey_item_id]
      scenarios = Survey::Scenario
        .where(survey_item_arel.eq(nil).or(survey_item_arel.in(user.survey_answers_survey_items.pluck(:id))))
        .joins(:triggers)
        .merge(
          Survey::Scenario::Trigger
            .where(event_name: event.class.name)
            .with_event_threshold_met(count)
        )

      scenarios.each do |scenario|
        maybe_new_member = scenario.members.find_or_initialize_by user: user
        maybe_new_member.save! if maybe_new_member.new_record?
      end

      # if event.class == LessonFinishedEvent
      #   binding.irb
      # end
    end
  end
end
