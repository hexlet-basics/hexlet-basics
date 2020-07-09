# frozen_string_literal: true

class Lesson::MemberService
  def self.ensure_member!(lesson:, lesson_version:, language:, user:)
    return Language::Lesson::Guest.new if user.guest?

    LessonMemberMutator.find_or_create_member!(
      lesson: lesson,
      lesson_version: lesson_version,
      user: user,
      language: language
    )
  end
end
