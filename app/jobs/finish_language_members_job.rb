class FinishLanguageMembersJob < ApplicationJob
  sidekiq_options retry: 0

  def perform(offset = 0, limit = 10_000)
    started_lesson_members = Language::Member.started
                                             .limit(limit)
                                             .offset(offset)
                                             .order(:id)

    started_lesson_members.find_each do |lesson_member|
      lesson_member.finish! if lesson_member.all_lessons_finished?
    end

    if started_lesson_members.any?
      offset += limit

      FinishLanguageMembersJob.set(wait: 1.minute).perform_later(offset)
    end
  end
end
