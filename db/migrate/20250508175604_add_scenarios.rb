class AddScenarios < ActiveRecord::Migration[8.0]
  def change
    scenarios_data = [
      {
        name: 'Входной сценарий',
        triggers: [
          { event_name: 'UserSignedUpEvent', event_threshold_count: 1 }
        ],
        survey_slugs: %w[
        goal
        coding-experience
      ]
      },
      {
        name: 'Смена карьеры',
        survey_item_slug: 'goal-item1',
        triggers: [
          { event_name: 'LessonFinishedEvent', event_threshold_count: 3 },
          { event_name: 'BookRequestedEvent', event_threshold_count: 1 }
        ],
        survey_slugs: %w[
        career-change-study-plan
        career-change-reason
        career-change-barrier
        career-change-priority
        career-change-preferred-intro-format
        career-change-contact-method
      ]
      },
      {
        name: 'Учеба',
        survey_item_slug: 'goal-item3',
        triggers: [
          { event_name: 'LessonFinishedEvent', event_threshold_count: 3 },
          { event_name: 'BookRequestedEvent', event_threshold_count: 1 }
        ],
        survey_slugs: %w[
        formal-study-institution
        formal-study-context
        formal-study-depth
        formal-study-priority
      ]
      },
      {
        name: 'Повышение квалификации',
        survey_item_slug: 'goal-item2',
        triggers: [
          { event_name: 'LessonFinishedEvent', event_threshold_count: 3 },
          { event_name: 'BookRequestedEvent', event_threshold_count: 1 }
        ],
        survey_slugs: %w[
        new-skill-current-role
        new-skill-task-type
        new-skill-depth
      ]
      }
    ]


    scenarios_data.each do |data|
      scenario = Survey::Scenario.find_or_initialize_by name: data[:name]
      scenario.locale = :ru
      survey_item = Survey::Item.find_by(slug: data[:survey_item_slug]) if data[:survey_item_slug]
      scenario.survey_item = survey_item if survey_item
      scenario.save!

      data[:triggers].each do |trigger_data|
        trigger = scenario.triggers.find_or_initialize_by(event_name: trigger_data[:event_name])
        trigger.event_threshold_count = trigger_data[:event_threshold_count]
        trigger.save!
      end

      data[:survey_slugs].each do |slug|
        survey = Survey.find_by! slug: slug
        item = scenario.items.find_or_initialize_by survey: survey
        item.save! if item.new_record?
      end
    end

    answers = Survey::Answer.includes(:user, survey: [ :scenarios ])
    answers.find_each do |answer|
      user = answer.user
      survey = answer.survey

      survey.scenarios.each do |scenario|
        member = scenario.members.find_or_initialize_by(user: user)
        member.save! if member.new_record?

        unless member.next_survey
          member.state = :finished
          member.save!
        end
      end
    end
  end
end
