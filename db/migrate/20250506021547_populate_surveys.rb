class PopulateSurveys < ActiveRecord::Migration[8.0]
  def change
    Survey.where(slug: %w[goal coding-experience]).update_all(run_always: true)

    Survey.where(slug: %w[
      career-change-reason
      career-change-study-plan
      new-skill-current-role
      new-skill-task-type
      new-skill-depth
      formal-study-institution
      formal-study-context
      formal-study-depth
      formal-study-priority
    ]).update_all(run_after_finishing_lessons_count: 3)

    Survey.where(slug: %w[
      career-change-barrier
      career-change-time-commitment
    ]).update_all(run_after_finishing_lessons_count: 6)

    Survey.where(slug: %w[
      career-change-priority
      preferred-intro-format
      career-change-contact-method
    ]).update_all(run_after_finishing_lessons_count: 9)
  end
end
