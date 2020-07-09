# frozen_string_literal: true

class LessonMemberMutator
  def self.find_or_create_member!(lesson:, lesson_version:, language:, user:)
    member = Language::Lesson::Member.find_by(language: language, user: user, lesson: lesson)
    return member if member.present?

    member = Language::Lesson::Member.new(
      lesson: lesson,
      lesson_version: lesson_version,
      user: user,
      language: language
    )
    member.save!

    member
  end
end
